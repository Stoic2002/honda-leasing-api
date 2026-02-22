package handler

import (
	"net/http"
	"strconv"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/officer"
	"honda-leasing-api/pkg/pagination"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type OfficerHandler struct {
	service officer.Service
}

func NewOfficerHandler(service officer.Service) *OfficerHandler {
	return &OfficerHandler{service: service}
}

func (h *OfficerHandler) GetIncomingOrders(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	pg := contract.PaginationFilter{
		Page:  page,
		Limit: limit,
	}

	orders, total, err := h.service.GetIncomingOrders(c.Request.Context(), pg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch incoming orders"))
		return
	}

	// Map entities to DTOs
	var orderResponses []IncomingOrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, toIncomingOrderResponse(o))
	}

	meta := pagination.BuildMeta(page, limit, total)
	c.JSON(http.StatusOK, response.SuccessPaginated(http.StatusOK, "Successfully fetched pending orders", orderResponses, meta))
}

func (h *OfficerHandler) GetMyTasks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	pg := contract.PaginationFilter{
		Page:  page,
		Limit: limit,
	}

	// Get user's role from JWT context
	roleVal, _ := c.Get("role")
	userRole := roleVal.(string)

	tasks, total, err := h.service.GetMyTasks(c.Request.Context(), userRole, pg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch tasks"))
		return
	}

	// Map entities to DTOs
	var taskResponses []OfficerTaskResponse
	for _, t := range tasks {
		taskResponses = append(taskResponses, toOfficerTaskResponse(t))
	}

	meta := pagination.BuildMeta(page, limit, total)
	c.JSON(http.StatusOK, response.SuccessPaginated(http.StatusOK, "Successfully fetched your assigned tasks", taskResponses, meta))
}

func (h *OfficerHandler) ProcessTask(c *gin.Context) {
	idStr := c.Param("taskId")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid task ID format"))
		return
	}

	// Get user's role from JWT context
	roleVal, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, "Role context missing"))
		return
	}
	userRole := roleVal.(string)

	var req ProcessTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	// Map handler DTO to service input
	input := officer.ProcessTaskInput{
		Notes: req.Notes,
	}

	err = h.service.ProcessOrderTask(c.Request.Context(), taskID, userRole, input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Error(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Task successfully processed and moved to next stage", nil))
}
