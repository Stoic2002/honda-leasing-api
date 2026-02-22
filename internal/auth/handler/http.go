package handler

import (
	"net/http"

	"honda-leasing-api/internal/auth"
	"honda-leasing-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service auth.Service
}

func NewAuthHandler(service auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	resp, err := h.service.Login(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Login successful", toLoginResponse(resp)))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	newAccess, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Token refreshed successfully", gin.H{
		"access_token": newAccess,
	}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Logged out successfully", nil))
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized access"))
		return
	}

	uid := userID.(int64)

	profile, err := h.service.GetProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched user profile", toUserProfileResponse(profile)))
}
