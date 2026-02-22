package handler

import "github.com/gin-gonic/gin"

func (h *CatalogHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	v1 := router.Group("/api/v1/catalog")
	v1.Use(authMiddleware)
	{
		v1.GET("/motors", h.GetMotors)
		v1.GET("/motors/:id", h.GetMotorByID)
		v1.GET("/leasing-products", h.GetLeasingProducts)
	}
}
