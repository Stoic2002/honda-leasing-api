package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"honda-leasing-api/internal/delivery"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/pkg/pagination"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	maxUploadSize = 5 << 20 // 5MB
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".pdf":  true,
}

type DeliveryHandler struct {
	service delivery.Service
}

func NewDeliveryHandler(service delivery.Service) *DeliveryHandler {
	return &DeliveryHandler{service: service}
}

func (h *DeliveryHandler) GetDeliveryOrders(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	pg := contract.PaginationFilter{
		Page:  page,
		Limit: limit,
	}

	roleVal, _ := c.Get("role")
	userRole := roleVal.(string)

	orders, total, err := h.service.GetDeliveryOrders(c.Request.Context(), userRole, pg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch delivery orders"))
		return
	}

	var orderResponses []DeliveryOrderResponse
	for _, o := range orders {
		orderResponses = append(orderResponses, toDeliveryOrderResponse(o))
	}

	meta := pagination.BuildMeta(page, limit, total)
	c.JSON(http.StatusOK, response.SuccessPaginated(http.StatusOK, "Successfully fetched delivery orders", orderResponses, meta))
}

func (h *DeliveryHandler) GetDeliveryTasks(c *gin.Context) {
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

	tasks, total, err := h.service.GetDeliveryTasks(c.Request.Context(), userRole, pg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch delivery tasks"))
		return
	}

	// Map entities to DTOs
	var taskResponses []DeliveryTaskResponse
	for _, t := range tasks {
		taskResponses = append(taskResponses, toDeliveryTaskResponse(t))
	}

	meta := pagination.BuildMeta(page, limit, total)

	c.JSON(http.StatusOK, response.SuccessPaginated(http.StatusOK, "Successfully fetched pending delivery tasks", taskResponses, meta))
}

func (h *DeliveryHandler) CompleteDeliveryTask(c *gin.Context) {
	idStr := c.Param("taskId")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid task ID format"))
		return
	}

	// Get user's role from JWT context
	roleVal, _ := c.Get("role")
	userRole := roleVal.(string)

	// Handle multipart form-data upload for file
	file, err := c.FormFile("foto_serah_terima")
	var req delivery.CompleteDeliveryRequest

	if err == nil && file != nil {
		// Validate file size
		if file.Size > maxUploadSize {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "File size exceeds maximum of 5MB"))
			return
		}

		// Validate file extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "File type not allowed. Allowed: jpg, jpeg, png, pdf"))
			return
		}

		// Sanitize filename — strip path components to prevent path traversal
		safeFilename := filepath.Base(file.Filename)
		safeFilename = strings.ReplaceAll(safeFilename, "..", "")

		// Create uploads dir if not exists
		uploadDir := "./uploads"
		if _, statErr := os.Stat(uploadDir); os.IsNotExist(statErr) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		filename := fmt.Sprintf("delivery_%d_%s", taskID, safeFilename)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to save file"))
			return
		}

		req = delivery.CompleteDeliveryRequest{
			FileName: filename,
			FileSize: float64(file.Size),
			FileType: ext[1:], // Remove the dot
			FileURL:  fmt.Sprintf("/uploads/%s", filename),
		}
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Error processing file upload: "+err.Error()))
		return
	}

	err = h.service.CompleteDelivery(c.Request.Context(), taskID, userRole, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Error(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Motor delivery task successfully completed", nil))
}
