package handlers

import (
	"net/http"
	"strconv"

	"emarket-api/middlewares"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/market/domains"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase domains.UseCase
}

func (h *Handler) GetGISConfig(c *gin.Context) {
	res, err := h.useCase.GetGISConfig(c.Param("id"), false)
	if err != nil {
		respondError(c, err)
		return
	}
	// Published GIS payloads are large but immutable until the matching PUT.
	// Let the browser reuse the compressed response on a refresh instead of
	// downloading and parsing the same multi-megabyte JSON again immediately.
	c.Header("Cache-Control", "private, max-age=300, stale-while-revalidate=60")
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

func (h *Handler) GetGISDraft(c *gin.Context) {
	res, err := h.useCase.GetGISConfig(c.Param("id"), true)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

func (h *Handler) SaveGISConfig(c *gin.Context) {
	var req requests.UpdateMarketGISRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}
	res, err := h.useCase.SaveGISConfig(c.Param("id"), middlewares.GetUserID(c), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

func (h *Handler) GetUtilitySettings(c *gin.Context) {
	res, err := h.useCase.GetUtilitySettings(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

func (h *Handler) SaveUtilitySettings(c *gin.Context) {
	var req requests.UpdateMarketUtilitySettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}
	res, err := h.useCase.SaveUtilitySettings(c.Param("id"), middlewares.GetUserID(c), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

func NewHandler(uc domains.UseCase) *Handler {
	return &Handler{useCase: uc}
}

func respondError(c *gin.Context, err error) {
	if appErr, ok := err.(*responses.AppErr); ok {
		c.JSON(appErr.StatusCode, responses.FailRespone(appErr))
		return
	}
	c.JSON(http.StatusBadRequest, responses.FailRespone(err))
}

// GetMarkets godoc
// @Summary ดึงรายการตลาดทั้งหมด
// @Tags Markets
// @ID get-markets
// @Produce json
// @Param search query string false "คำค้นหา (ชื่อตลาด/เขต/รหัส)"
// @Success 200 {object} responses.Success{data=[]models.Market}
// @Router /markets [get]
func (h *Handler) GetMarkets(c *gin.Context) {
	search := c.Query("search")
	res, err := h.useCase.GetMarkets(search)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

func (h *Handler) GetMyMarkets(c *gin.Context) {
	res, err := h.useCase.GetMarketsForUser(c.Query("search"), middlewares.GetUserID(c), middlewares.GetUserRole(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// GetMarketByID godoc
// @Summary ดึงข้อมูลตลาดรายตัวตาม ID หรือ Slug
// @Tags Markets
// @ID get-market-by-id
// @Produce json
// @Param id path string true "Market ID หรือ Slug"
// @Success 200 {object} responses.Success{data=models.Market}
// @Router /markets/{id} [get]
func (h *Handler) GetMarketByID(c *gin.Context) {
	idOrSlug := c.Param("id")
	res, err := h.useCase.GetMarketByIDOrSlug(idOrSlug)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// GetMarketZones godoc
// @Summary ดึงรายการโซนของตลาด
// @Tags Markets
// @ID get-market-zones
// @Produce json
// @Param id path int true "Market ID"
// @Success 200 {object} responses.Success{data=[]models.Zone}
// @Router /markets/{id}/zones [get]
func (h *Handler) GetMarketZones(c *gin.Context) {
	idStr := c.Param("id")
	marketID := 0
	if idStr != "" {
		var err error
		marketID, err = strconv.Atoi(idStr)
		if err != nil {
			market, mErr := h.useCase.GetMarketByIDOrSlug(idStr)
			if mErr != nil {
				respondError(c, responses.NewAppErr(http.StatusBadRequest, "Market ID ไม่ถูกต้อง"))
				return
			}
			marketID = market.ID
		}
	}

	res, err := h.useCase.GetMarketZones(marketID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}
