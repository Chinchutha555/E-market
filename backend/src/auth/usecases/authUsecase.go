package usecases

import (
	"errors"
	"net/http"

	"emarket-api/helpers"
	"emarket-api/logs"
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/auth/domains"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type useCase struct {
	repo domains.Repository
}

func NewUseCase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) toUserResponse(user *models.User) responses.UserResponse {
	res := responses.UserResponse{
		ID: user.ID, Username: user.Username, FullName: user.FullName,
		Role: user.Role, Email: user.Email, Phone: user.Phone,
		Status: user.Status, CreatedAt: user.CreatedAt,
	}
	if user.Role == "staff" || user.Role == "finance" || user.Role == "market_director" {
		if markets, err := u.repo.FindAssignedMarkets(user.ID); err == nil {
			for _, market := range markets {
				res.AssignedMarkets = append(res.AssignedMarkets, responses.AssignedMarketResponse{
					ID: market.ID, Slug: market.Slug, Code: market.Code, Name: market.Name, District: market.District,
				})
			}
		}
	}
	return res
}

func (u *useCase) Login(req requests.LoginRequest) (*responses.LoginResponse, error) {
	user, err := u.repo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusUnauthorized, "ชื่อผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil && user.Password != req.Password {
		return nil, responses.NewAppErr(http.StatusUnauthorized, "ชื่อผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง")
	}

	token, err := helpers.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		refreshToken = token
	}

	res := responses.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		Role:         user.Role,
		User:         u.toUserResponse(user),
	}
	return &res, nil
}

func (u *useCase) GetProfile(userID int) (*responses.UserResponse, error) {
	user, err := u.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลผู้ใช้งาน")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	res := u.toUserResponse(user)
	return &res, nil
}
