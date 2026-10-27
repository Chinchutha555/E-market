package handlers

import (
	"net/http"
	"strconv"

	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/application_type/domains"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase domains.UseCase
}

func NewHandler(uc domains.UseCase) *Handler {
	return &Handler{useCase: uc}
}

func respondError(c *gin.Context, err error) {
	if appErr, ok := err.(*responses.AppErr); ok {
		if len(appErr.Fields) > 0 {
			c.JSON(appErr.StatusCode, responses.ValidateResponse(appErr.Fields))
			return
		}
		c.JSON(appErr.StatusCode, responses.FailRespone(appErr))
		return
	}
	c.JSON(http.StatusBadRequest, responses.FailRespone(err))
}

// GetApplicationTypes godoc
// @Summary ดึงรายการประเภทคำร้องทั้งหมด
// @Tags ApplicationTypes
// @ID get-application-types
// @Produce json
// @Param active_only query bool false "Active only"
// @Param search query string false "Search keyword"
// @Success 200 {object} responses.Success{data=[]responses.ApplicationTypeResponse}
// @Router /application-types [get]
func (h *Handler) GetApplicationTypes(c *gin.Context) {
	var filter requests.ApplicationTypeFilter
	_ = c.ShouldBindQuery(&filter)

	res, err := h.useCase.GetApplicationTypes(filter)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// GetApplicationTypeByID godoc
// @Summary ดึงข้อมูลประเภทคำร้องตาม ID
// @Tags ApplicationTypes
// @ID get-application-type-by-id
// @Produce json
// @Param id path int true "Application Type ID"
// @Success 200 {object} responses.Success{data=responses.ApplicationTypeResponse}
// @Router /application-types/{id} [get]
func (h *Handler) GetApplicationTypeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	res, err := h.useCase.GetApplicationTypeByID(uint(id))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// CreateApplicationType godoc
// @Summary สร้างประเภทคำร้องใหม่ (Admin Only)
// @Tags ApplicationTypes
// @ID create-application-type
// @Accept json
// @Produce json
// @Param body body requests.CreateApplicationTypeRequest true "ข้อมูลประเภทคำร้อง"
// @Success 201 {object} responses.Success{data=responses.ApplicationTypeResponse}
// @Router /application-types [post]
func (h *Handler) CreateApplicationType(c *gin.Context) {
	var req requests.CreateApplicationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "รูปแบบข้อมูลไม่ถูกต้อง"))
		return
	}

	res, err := h.useCase.CreateApplicationType(req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(res, http.StatusCreated))
}

// UpdateApplicationType godoc
// @Summary แก้ไขประเภทคำร้อง (Admin Only)
// @Tags ApplicationTypes
// @ID update-application-type
// @Accept json
// @Produce json
// @Param id path int true "Application Type ID"
// @Param body body requests.UpdateApplicationTypeRequest true "ข้อมูลแก้ไข"
// @Success 200 {object} responses.Success{data=responses.ApplicationTypeResponse}
// @Router /application-types/{id} [put]
func (h *Handler) UpdateApplicationType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	var req requests.UpdateApplicationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "รูปแบบข้อมูลไม่ถูกต้อง"))
		return
	}

	res, err := h.useCase.UpdateApplicationType(uint(id), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// DeleteApplicationType godoc
// @Summary ลบประเภทคำร้อง (Admin Only)
// @Tags ApplicationTypes
// @ID delete-application-type
// @Produce json
// @Param id path int true "Application Type ID"
// @Success 200 {object} responses.Success
// @Router /application-types/{id} [delete]
func (h *Handler) DeleteApplicationType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	if err := h.useCase.DeleteApplicationType(uint(id)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"deleted": true}, http.StatusOK))
}

// ToggleActive godoc
// @Summary สลับสถานะเปิด-ปิดใช้งานประเภทคำร้อง (Admin Only)
// @Tags ApplicationTypes
// @ID toggle-application-type-active
// @Produce json
// @Param id path int true "Application Type ID"
// @Success 200 {object} responses.Success{data=responses.ApplicationTypeResponse}
// @Router /application-types/{id}/toggle [patch]
func (h *Handler) ToggleActive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	res, err := h.useCase.ToggleActive(uint(id))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}
