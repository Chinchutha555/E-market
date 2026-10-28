package handlers

import (
	"net/http"
	"strconv"

	"emarket-api/helpers"
	"emarket-api/middlewares"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/contract/domains"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
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

// GetContracts godoc
// @Summary ดึงรายการสัญญาเช่าทั้งหมด
// @Tags Contracts
// @ID get-contracts
// @Produce json
// @Param stall_ident query string false "Stall ident"
// @Param status query string false "Status"
// @Param search query string false "Search keyword"
// @Success 200 {object} responses.Success{data=[]responses.ContractResponse}
// @Router /contracts [get]
func (h *Handler) GetContracts(c *gin.Context) {
	var filter requests.ContractFilter
	_ = c.ShouldBindQuery(&filter)

	res, err := h.useCase.GetContracts(filter)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// GetContractByID godoc
// @Summary ดูรายละเอียดสัญญาตาม ID
// @Tags Contracts
// @ID get-contract-by-id
// @Produce json
// @Param id path int true "Contract ID"
// @Success 200 {object} responses.Success{data=responses.ContractResponse}
// @Router /contracts/{id} [get]
func (h *Handler) GetContractByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	res, err := h.useCase.GetContractByID(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// CreateContract godoc
// @Summary จัดทำสัญญาเช่าใหม่ (Admin / Staff)
// @Tags Contracts
// @ID create-contract
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body requests.CreateContractRequest true "ข้อมูลสัญญา"
// @Success 201 {object} responses.Success{data=responses.ContractResponse}
// @Router /contracts [post]
func (h *Handler) CreateContract(c *gin.Context) {
	userID := middlewares.GetUserID(c)
	userRole := middlewares.GetUserRole(c)
	var req requests.CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}
	if err := validator.Validate(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.ValidateResponse(helpers.ConverstError(err)))
		return
	}

	res, err := h.useCase.CreateContract(userID, userRole, req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(res, http.StatusCreated))
}

// ApproveContract godoc
// @Summary อนุมัติ/เปิดใช้งานสัญญาเช่า (Admin / Staff / Market Director)
// @Tags Contracts
// @ID approve-contract
// @Security BearerAuth
// @Produce json
// @Param id path int true "Contract ID"
// @Success 200 {object} responses.Success{data=responses.ContractResponse}
// @Router /contracts/{id}/approve [post]
func (h *Handler) ApproveContract(c *gin.Context) {
	userID := middlewares.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	res, err := h.useCase.ApproveContract(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(res, http.StatusOK))
}

// RejectContract godoc
// @Summary ปฏิเสธ/ส่งกลับสัญญาเช่า (Admin / Staff / Market Director)
// @Tags Contracts
// @ID reject-contract
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Contract ID"
// @Param body body requests.ReviewContractRequest false "เหตุผล"
// @Success 200 {object} responses.Success
// @Router /contracts/{id}/reject [post]
func (h *Handler) RejectContract(c *gin.Context) {
	userID := middlewares.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	var req requests.ReviewContractRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.useCase.RejectContract(id, userID, req.Reason); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"rejected": true}, http.StatusOK))
}

// TerminateContract godoc
// @Summary ยกเลิก/สิ้นสุดสัญญาเช่า
// @Tags Contracts
// @ID terminate-contract
// @Security BearerAuth
// @Produce json
// @Param id path int true "Contract ID"
// @Success 200 {object} responses.Success
// @Router /contracts/{id}/terminate [post]
func (h *Handler) TerminateContract(c *gin.Context) {
	userID := middlewares.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง"))
		return
	}

	if err := h.useCase.TerminateContract(id, userID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"terminated": true}, http.StatusOK))
}
