package handler

import (
	"net/http"
	"strconv"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/leasing"
	"honda-leasing-api/pkg/pagination"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type LeasingHandler struct {
	service leasing.Service
}

func NewLeasingHandler(service leasing.Service) *LeasingHandler {
	return &LeasingHandler{service: service}
}

func (h *LeasingHandler) SubmitOrder(c *gin.Context) {
	// Ambil userID dari JWT yang sudah di-set oleh middleware auth
	userIDVal, exists := c.Get("userID")
	if !exists {
		_ = c.Error(domain.ErrUnauthorized)
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		_ = c.Error(domain.ErrInternalServerError)
		return
	}

	var req SubmitOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	// Map handler DTO to service input (tanpa customer_id — diambil otomatis dari DB via userID)
	input := leasing.SubmitOrderInput{
		UserID:         userID,
		MotorID:        req.MotorID,
		ProductID:      req.ProductID,
		NilaiKendaraan: req.NilaiKendaraan,
		DpDibayar:      req.DpDibayar,
		TenorBulan:     req.TenorBulan,
	}

	cont, err := h.service.SubmitOrder(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	contractResp := toContractResponse(*cont)
	c.JSON(http.StatusCreated, response.Success(http.StatusCreated, "Order submitted successfully", contractResp))
}

func (h *LeasingHandler) GetMyOrders(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		_ = c.Error(domain.ErrUnauthorized)
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		_ = c.Error(domain.ErrInternalServerError)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	pg := contract.PaginationFilter{
		Page:  page,
		Limit: limit,
	}

	orders, total, err := h.service.GetMyOrders(c.Request.Context(), userID, pg)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var orderResponses []MyOrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, toMyOrderResponse(o))
	}

	meta := pagination.BuildMeta(page, limit, total)

	c.JSON(http.StatusOK, response.SuccessPaginated(http.StatusOK, "Successfully fetched your orders", orderResponses, meta))
}

func (h *LeasingHandler) GetContractProgress(c *gin.Context) {
	idStr := c.Param("id")
	contractID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	tasks, err := h.service.GetContractProgress(c.Request.Context(), contractID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Map entities to DTOs
	var taskResponses []TaskProgressResponse
	for _, t := range tasks {
		taskResponses = append(taskResponses, toTaskProgressResponse(t))
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched contract progress", taskResponses))
}
