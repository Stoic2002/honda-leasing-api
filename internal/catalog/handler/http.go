package handler

import (
	"net/http"
	"strconv"

	"honda-leasing-api/internal/catalog"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/pkg/pagination"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	service catalog.Service
}

func NewCatalogHandler(service catalog.Service) *CatalogHandler {
	return &CatalogHandler{service: service}
}

func (h *CatalogHandler) GetMotors(c *gin.Context) {
	search := c.Query("search")
	motorType := c.Query("type")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	filter := contract.CatalogFilter{
		Search:    search,
		MotorType: motorType,
	}

	pg := contract.PaginationFilter{
		Page:  page,
		Limit: limit,
	}

	motors, total, err := h.service.GetMotors(c.Request.Context(), filter, pg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch motors"))
		return
	}

	// Map entities to DTOs
	var motorResponses []MotorResponse
	for _, m := range motors {
		motorResponses = append(motorResponses, toMotorResponse(m))
	}

	meta := pagination.BuildMeta(page, limit, total)

	c.JSON(http.StatusOK, response.SuccessPaginated(http.StatusOK, "Successfully fetched motors", motorResponses, meta))
}

func (h *CatalogHandler) GetMotorByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid motor ID"))
		return
	}

	motor, err := h.service.GetMotorByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Motor not found"))
		return
	}

	motorResp := toMotorResponse(*motor)
	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched motor", motorResp))
}

func (h *CatalogHandler) GetLeasingProducts(c *gin.Context) {
	products, err := h.service.GetLeasingProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to fetch leasing products"))
		return
	}

	// Map entities to DTOs
	var productResponses []LeasingProductResponse
	for _, p := range products {
		productResponses = append(productResponses, toLeasingProductResponse(p))
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched leasing products", productResponses))
}
