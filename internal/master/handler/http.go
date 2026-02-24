package handler

import (
	"net/http"
	"strconv"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/internal/master"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	service master.Service
}

func NewMasterHandler(service master.Service) *MasterHandler {
	return &MasterHandler{service: service}
}

func (h *MasterHandler) GetProvinces(c *gin.Context) {
	provinces, err := h.service.GetProvinces(c.Request.Context())
	if err != nil {
		_ = c.Error(domain.ErrInternalServerError)
		return
	}

	var res []ProvinceResponse
	for _, p := range provinces {
		res = append(res, ProvinceResponse{
			ProvID:   p.ProvID,
			ProvName: p.ProvName,
		})
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched provinces", res))
}

func (h *MasterHandler) GetKabupatens(c *gin.Context) {
	provIDStr := c.Query("prov_id")
	if provIDStr == "" {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	provID, err := strconv.ParseInt(provIDStr, 10, 64)
	if err != nil {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	kabs, err := h.service.GetKabupatens(c.Request.Context(), provID)
	if err != nil {
		_ = c.Error(domain.ErrInternalServerError)
		return
	}

	var res []KabupatenResponse
	for _, k := range kabs {
		res = append(res, KabupatenResponse{
			KabID:   k.KabID,
			KabName: k.KabName,
		})
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched kabupatens", res))
}

func (h *MasterHandler) GetKecamatans(c *gin.Context) {
	kabIDStr := c.Query("kab_id")
	if kabIDStr == "" {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	kabID, err := strconv.ParseInt(kabIDStr, 10, 64)
	if err != nil {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	kecs, err := h.service.GetKecamatans(c.Request.Context(), kabID)
	if err != nil {
		_ = c.Error(domain.ErrInternalServerError)
		return
	}

	var res []KecamatanResponse
	for _, k := range kecs {
		res = append(res, KecamatanResponse{
			KecID:   k.KecID,
			KecName: k.KecName,
		})
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched kecamatans", res))
}

func (h *MasterHandler) GetKelurahans(c *gin.Context) {
	kecIDStr := c.Query("kec_id")
	if kecIDStr == "" {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	kecID, err := strconv.ParseInt(kecIDStr, 10, 64)
	if err != nil {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	kels, err := h.service.GetKelurahans(c.Request.Context(), kecID)
	if err != nil {
		_ = c.Error(domain.ErrInternalServerError)
		return
	}

	var res []KelurahanResponse
	for _, k := range kels {
		res = append(res, KelurahanResponse{
			KelID:   k.KelID,
			KelName: k.KelName,
		})
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched kelurahans", res))
}
