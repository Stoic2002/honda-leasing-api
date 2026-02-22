package crypto

import (
	"fmt"
	"time"

	"honda-leasing-api/configs"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserID   int64  `json:"userId"`
	Email    string `json:"email"`
	RoleName string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID int64, email, role string, cfg configs.JwtConfig) (accessToken string, refreshToken string, err error) {
	// Access Token
	accessExp := time.Now().Add(time.Minute * time.Duration(cfg.ExpireMinutes))
	accessClaims := &JwtCustomClaims{
		UserID:   userID,
		Email:    email,
		RoleName: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshExp := time.Now().Add(time.Hour * 24 * time.Duration(cfg.RefreshDays))
	refreshClaims := &JwtCustomClaims{
		UserID: userID,
		// We don't embed strictly all data in refresh usually, but keeping some identity is fine
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString, secret string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
