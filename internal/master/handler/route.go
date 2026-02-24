package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *MasterHandler) RegisterRoutes(router *gin.Engine) {
	// Master data is publicly accessible for frontend dropdowns
	v1 := router.Group("/api/v1/master")
	{
		v1.GET("/provinces", h.GetProvinces)
		v1.GET("/kabupatens", h.GetKabupatens)
		v1.GET("/kecamatans", h.GetKecamatans)
		v1.GET("/kelurahans", h.GetKelurahans)
	}
}
