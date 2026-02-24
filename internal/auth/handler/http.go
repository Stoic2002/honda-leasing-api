package handler

import (
	"net/http"

	"honda-leasing-api/internal/auth"
	"honda-leasing-api/internal/domain"
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
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	resp, err := h.service.Login(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Login successful", toLoginResponse(resp)))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(domain.ErrInvalidInput)
		return
	}

	newAccess, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
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
		_ = c.Error(domain.ErrUnauthorized)
		return
	}

	uid := userID.(int64)

	profile, err := h.service.GetProfile(c.Request.Context(), uid)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Success(http.StatusOK, "Successfully fetched user profile", toUserProfileResponse(profile)))
}
