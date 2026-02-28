package handler

import (
	"net/http"
	"strconv"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/internal/finance"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	service finance.Service
}

func NewFinanceHandler(service finance.Service) *FinanceHandler {
	return &FinanceHandler{service: service}
}

func (h *FinanceHandler) GetSchedules(c *gin.Context) {
	contractIDStr := c.Query("contract_id")
	var contractID int64
	var err error

	if contractIDStr != "" {
		contractID, err = strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			_ = c.Error(domain.ErrInvalidInput)
			return
		}
	}

	res, err := h.service.GetPaymentSchedules(c.Request.Context(), contractID)
	if err != nil {
		_ = c.Error(domain.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched payment schedules", res))
}
