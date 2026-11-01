package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"emarket-api/logs"
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/market/domains"
	"gorm.io/gorm"
)

func resolveMarketJSONPath(slug string) string {
	if slug == "" {
		slug = "chatuchak"
	}
	candidates := []string{
		filepath.Join("data", "markets", slug+".json"),
		filepath.Join("backend", "data", "markets", slug+".json"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return candidates[0]
}

type baseProjectsCacheEntry struct {
	modTime  time.Time
	size     int64
	projects json.RawMessage
}

var (
	baseProjectsCache   = map[string]baseProjectsCacheEntry{}
	baseProjectsCacheMu sync.Mutex
)

// อ่านผังโครงการจากไฟล์ JSON โดยแคชผลที่แปลงแล้ว จนกว่าไฟล์จะถูกแก้ไข (ไฟล์ ~1MB/ตลาด อ่าน+แปลงทุกคำขอช้า)
func loadBaseProjectsFromDisk(slug string) json.RawMessage {
	jsonPath := resolveMarketJSONPath(slug)
	info, statErr := os.Stat(jsonPath)
	if statErr == nil {
		baseProjectsCacheMu.Lock()
		entry, ok := baseProjectsCache[jsonPath]
		baseProjectsCacheMu.Unlock()
		if ok && entry.modTime.Equal(info.ModTime()) && entry.size == info.Size() {
			return entry.projects
		}
	}
	rawBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return json.RawMessage("[]")
	}
	var marketDoc map[string]json.RawMessage
	if err := json.Unmarshal(rawBytes, &marketDoc); err != nil {
		return json.RawMessage("[]")
	}
	projects := json.RawMessage("[]")
	if bpRaw, ok := marketDoc["projects"]; ok && json.Valid(bpRaw) {
		var compact bytes.Buffer
		if err := json.Compact(&compact, bpRaw); err == nil {
			projects = json.RawMessage(compact.Bytes())
		} else {
			projects = bpRaw
		}
	}
	if statErr == nil {
		baseProjectsCacheMu.Lock()
		baseProjectsCache[jsonPath] = baseProjectsCacheEntry{modTime: info.ModTime(), size: info.Size(), projects: projects}
		baseProjectsCacheMu.Unlock()
	}
	return projects
}

func mergeProjectsWithBase(baseProjectsRaw json.RawMessage, incomingProjectsRaw json.RawMessage) json.RawMessage {
	var incomingProjects []map[string]interface{}
	if err := json.Unmarshal(incomingProjectsRaw, &incomingProjects); err != nil || len(incomingProjects) == 0 {
		if len(baseProjectsRaw) > 2 {
			return baseProjectsRaw
		}
		return json.RawMessage("[]")
	}

	var baseProjects []map[string]interface{}
	_ = json.Unmarshal(baseProjectsRaw, &baseProjects)
	if len(baseProjects) == 0 {
		return incomingProjectsRaw
	}

	baseIDs := make(map[string]bool)
	for _, bp := range baseProjects {
		if id, ok := bp["id"].(string); ok && !strings.HasPrefix(id, "zone-") {
			baseIDs[id] = true
		}
	}

	incomingHasBase := false
	for _, ip := range incomingProjects {
		if id, ok := ip["id"].(string); ok && baseIDs[id] {
			incomingHasBase = true
			break
		}
	}

	mergedProjects := incomingProjects
	if !incomingHasBase && len(baseProjects) > 0 {
		existingNonZone := make([]map[string]interface{}, 0, len(baseProjects))
		for _, bp := range baseProjects {
			if id, ok := bp["id"].(string); ok && !strings.HasPrefix(id, "zone-") {
				existingNonZone = append(existingNonZone, bp)
			}
		}
		mergedProjects = append(existingNonZone, incomingProjects...)
	}

	if mergedBytes, err := json.Marshal(mergedProjects); err == nil {
		return json.RawMessage(mergedBytes)
	}
	return incomingProjectsRaw
}

// syncMarketJSONFile merges incoming projects with the base JSON file in
// backend/data/markets/<slug>.json and writes the unified file back to disk
// so all existing and newly drawn zones/stalls live in the same JSON file.
func syncMarketJSONFile(slug string, incomingProjectsRaw json.RawMessage) json.RawMessage {
	jsonPath := resolveMarketJSONPath(slug)
	baseProjectsRaw := loadBaseProjectsFromDisk(slug)
	mergedRaw := mergeProjectsWithBase(baseProjectsRaw, incomingProjectsRaw)

	rawBytes, err := os.ReadFile(jsonPath)
	if err == nil {
		var marketDoc map[string]interface{}
		if json.Unmarshal(rawBytes, &marketDoc) == nil {
			var mergedProjects []map[string]interface{}
			if json.Unmarshal(mergedRaw, &mergedProjects) == nil {
				marketDoc["projects"] = mergedProjects
				if outBytes, err := json.MarshalIndent(marketDoc, "", "  "); err == nil {
					_ = os.WriteFile(jsonPath, outBytes, 0644)
				}
			}
		}
	}
	return mergedRaw
}

type useCase struct {
	repo domains.Repository
}

func gisResponse(config *models.MarketGISConfig) *responses.MarketGISResponse {
	projects, facilities := json.RawMessage(config.ProjectsJSON), json.RawMessage(config.FacilitiesJSON)
	if !json.Valid(projects) {
		projects = json.RawMessage("[]")
	}
	if !json.Valid(facilities) {
		facilities = json.RawMessage("[]")
	}
	return &responses.MarketGISResponse{MarketID: config.MarketID, Projects: projects, Facilities: facilities, Published: config.Published, UpdatedBy: config.UpdatedBy, UpdatedAt: config.UpdatedAt}
}

func NewUseCase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) GetMarkets(search string) ([]models.Market, error) {
	markets, err := u.repo.FindAll(search)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return markets, nil
}

func (u *useCase) GetMarketsForUser(search string, userID int, role string) ([]models.Market, error) {
	if role == "admin" {
		return u.GetMarkets(search)
	}
	if role != "staff" {
		return nil, responses.NewAppErr(http.StatusForbidden, "ไม่มีสิทธิ์เข้าถึงข้อมูลตลาดสำหรับเจ้าหน้าที่")
	}
	markets, err := u.repo.FindByStaffUserID(search, userID)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return markets, nil
}

func (u *useCase) GetMarketByIDOrSlug(idOrSlug string) (*models.Market, error) {
	market, err := u.repo.FindByIDOrSlug(idOrSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลตลาด")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return market, nil
}

func (u *useCase) GetMarketZones(marketID int) ([]models.Zone, error) {
	zones, err := u.repo.FindZonesByMarketID(marketID)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return zones, nil
}

type gisMetaReader interface {
	FindGISConfigMeta(marketID int) (time.Time, bool, error)
}

type gisResponseCacheEntry struct {
	updatedAt time.Time
	baseSig   string
	resp      *responses.MarketGISResponse
}

var (
	gisResponseCache   = map[string]gisResponseCacheEntry{}
	gisResponseCacheMu sync.Mutex
)

// The merged GIS response is several MB; merging it (unmarshal + marshal) on every
// request took 1-2s. Reuse it until the stored config or the base JSON file changes.
func (u *useCase) GetGISConfig(idOrSlug string, includeDraft bool) (*responses.MarketGISResponse, error) {
	market, err := u.repo.FindByIDOrSlug(idOrSlug)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลตลาด")
	}
	baseProjects := loadBaseProjectsFromDisk(market.Slug)
	reader, canCache := u.repo.(gisMetaReader)
	if !canCache {
		return u.buildGISConfig(market, baseProjects, includeDraft)
	}
	updatedAt, _, metaErr := reader.FindGISConfigMeta(market.ID)
	if metaErr != nil {
		return u.buildGISConfig(market, baseProjects, includeDraft)
	}
	key := fmt.Sprintf("%d:%t", market.ID, includeDraft)
	baseSig := fmt.Sprintf("%p:%d", baseProjects, len(baseProjects))
	gisResponseCacheMu.Lock()
	entry, ok := gisResponseCache[key]
	gisResponseCacheMu.Unlock()
	if ok && entry.updatedAt.Equal(updatedAt) && entry.baseSig == baseSig {
		return entry.resp, nil
	}
	resp, err := u.buildGISConfig(market, baseProjects, includeDraft)
	if err != nil {
		return nil, err
	}
	gisResponseCacheMu.Lock()
	gisResponseCache[key] = gisResponseCacheEntry{updatedAt: updatedAt, baseSig: baseSig, resp: resp}
	gisResponseCacheMu.Unlock()
	return resp, nil
}

func (u *useCase) buildGISConfig(market *models.Market, baseProjects json.RawMessage, includeDraft bool) (*responses.MarketGISResponse, error) {
	config, err := u.repo.FindGISConfig(market.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &responses.MarketGISResponse{
			MarketID:   market.ID,
			Projects:   baseProjects,
			Facilities: json.RawMessage("[]"),
			Published:  true,
		}, nil
	}
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	if !includeDraft && !config.Published {
		return &responses.MarketGISResponse{
			MarketID:   market.ID,
			Projects:   baseProjects,
			Facilities: json.RawMessage("[]"),
			Published:  true,
		}, nil
	}
	resp := gisResponse(config)
	resp.Projects = mergeProjectsWithBase(baseProjects, resp.Projects)
	return resp, nil
}

func (u *useCase) SaveGISConfig(idOrSlug string, userID int, req requests.UpdateMarketGISRequest) (*responses.MarketGISResponse, error) {
	if !json.Valid(req.Projects) || !json.Valid(req.Facilities) {
		return nil, responses.NewAppErr(http.StatusBadRequest, "ข้อมูล GIS ไม่ถูกต้อง")
	}
	market, err := u.repo.FindByIDOrSlug(idOrSlug)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลตลาด")
	}
	mergedProjects := syncMarketJSONFile(market.Slug, req.Projects)
	config := &models.MarketGISConfig{MarketID: market.ID, ProjectsJSON: string(mergedProjects), FacilitiesJSON: string(req.Facilities), Published: req.Published, CreatedBy: &userID, UpdatedBy: &userID}
	if err := u.repo.SaveGISConfig(config); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return gisResponse(config), nil
}

// ตลาดขนาดใหญ่สามารถมีมิเตอร์หลายพันรายการในเอกสารเดียว จึงต้องเผื่อ payload
// มากกว่า 2MB เดิม ไม่เช่นนั้นการเพิ่มมิเตอร์แผงเดียวอาจถูกปฏิเสธทั้งเอกสาร
const maxUtilitySettingsBytes = 10 << 20

func (u *useCase) GetUtilitySettings(idOrSlug string) (*responses.MarketUtilitySettingsResponse, error) {
	market, err := u.repo.FindByIDOrSlug(idOrSlug)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลตลาด")
	}
	setting, err := u.repo.FindUtilitySettings(market.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ยังไม่เคยตั้งค่า: ส่ง settings เป็น null ให้ฝั่งหน้าจอใช้ค่าเริ่มต้น
			return &responses.MarketUtilitySettingsResponse{MarketID: market.ID, Settings: json.RawMessage("null")}, nil
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	settings := json.RawMessage(setting.SettingsJSON)
	if !json.Valid(settings) {
		settings = json.RawMessage("null")
	}
	updatedAt := setting.UpdatedAt
	return &responses.MarketUtilitySettingsResponse{MarketID: market.ID, Settings: settings, UpdatedBy: setting.UpdatedBy, UpdatedAt: &updatedAt}, nil
}

func (u *useCase) SaveUtilitySettings(idOrSlug string, userID int, req requests.UpdateMarketUtilitySettingsRequest) (*responses.MarketUtilitySettingsResponse, error) {
	if !json.Valid(req.Settings) || len(req.Settings) > maxUtilitySettingsBytes {
		return nil, responses.NewAppErr(http.StatusBadRequest, "ข้อมูลการตั้งค่าอัตราไม่ถูกต้อง")
	}
	market, err := u.repo.FindByIDOrSlug(idOrSlug)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลตลาด")
	}
	setting := &models.MarketUtilitySetting{MarketID: market.ID, SettingsJSON: string(req.Settings), UpdatedBy: &userID}
	if err := u.repo.SaveUtilitySettings(setting); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	now := time.Now()
	return &responses.MarketUtilitySettingsResponse{MarketID: market.ID, Settings: req.Settings, UpdatedBy: &userID, UpdatedAt: &now}, nil
}
