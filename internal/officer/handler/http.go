package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"honda-leasing-api/internal/domain"
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
		_ = c.Error(err)
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
		_ = c.Error(err)
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
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	// Get user's role from JWT context
	roleVal, exists := c.Get("role")
	if !exists {
		_ = c.Error(domain.ErrUnauthorized)
		return
	}
	userRole := roleVal.(string)

	// Instead of JSON bind, we read from Multipart Form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit Max
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	notes := c.PostForm("notes")
	dynamicAttributes := make(map[string]string)

	// 1. Process standard text attributes from form (e.g., attributes[Keterangan Wawancara])
	for key, values := range c.Request.MultipartForm.Value {
		// Key matching pattern like: attributes[NamaAtribut]
		if len(key) > 11 && key[:11] == "attributes[" && key[len(key)-1] == ']' {
			attributeName := key[11 : len(key)-1]
			if len(values) > 0 {
				dynamicAttributes[attributeName] = values[0]
			}
		}
	}

	// 2. Process file uploads from form
	for key, fileHeaders := range c.Request.MultipartForm.File {
		if len(key) > 11 && key[:11] == "attributes[" && key[len(key)-1] == ']' {
			attributeName := key[11 : len(key)-1]

			if len(fileHeaders) > 0 {
				fileHeader := fileHeaders[0]

				// Construct unique path
				fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
				uploadPath := filepath.Join("uploads", fileName)

				if err := c.SaveUploadedFile(fileHeader, uploadPath); err != nil {
					_ = c.Error(domain.ErrInternalServerError)
					return
				}

				// Generate accessible URL route
				fileUrl := fmt.Sprintf("/uploads/%s", fileName)
				dynamicAttributes[attributeName] = fileUrl
			}
		}
	}

	// Map handler DTO to service input
	input := officer.ProcessTaskInput{
		Notes:      notes,
		Attributes: dynamicAttributes,
	}

	err = h.service.ProcessOrderTask(c.Request.Context(), taskID, userRole, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Task successfully processed and moved to next stage", nil))
}
