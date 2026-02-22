package auth

import (
	"context"
	"fmt"

	"honda-leasing-api/configs"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/pkg/crypto"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	Role         string
}

type UserProfile struct {
	UserID   int64
	Email    string
	FullName string
	Role     string
}

type Service interface {
	Login(ctx context.Context, req LoginInput) (*LoginResult, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	GetProfile(ctx context.Context, userID int64) (*UserProfile, error)
}

type service struct {
	repo contract.AuthRepository
	cfg  configs.JwtConfig
}

func NewService(repo contract.AuthRepository, cfg configs.JwtConfig) Service {
	return &service{repo: repo, cfg: cfg}
}

func (s *service) Login(ctx context.Context, req LoginInput) (*LoginResult, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !crypto.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is disabled")
	}

	var roleName string
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].RoleName
	}
	if roleName == "" {
		roleName = "CUSTOMER" // Fallback safety
	}

	acc, ref, err := crypto.GenerateTokens(user.UserID, user.Email, roleName, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &LoginResult{
		AccessToken:  acc,
		RefreshToken: ref,
		Role:         roleName,
	}, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	// 1. Validate the old refresh token
	claims, err := crypto.ValidateToken(refreshToken, s.cfg.Secret)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// 2. Fetch the latest user info implicitly checking if they are still active
	user, err := s.repo.FindUserByID(ctx, claims.UserID)
	if err != nil || !user.IsActive {
		return "", fmt.Errorf("user no longer active or exists")
	}

	var roleName string
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].RoleName
	}

	// 3. Generate a new set of tokens (we only return the new access token)
	acc, _, err := crypto.GenerateTokens(user.UserID, user.Email, roleName, s.cfg)
	if err != nil {
		return "", err
	}

	return acc, nil
}

func (s *service) GetProfile(ctx context.Context, userID int64) (*UserProfile, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	var roleName string
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].RoleName
	}

	return &UserProfile{
		UserID:   user.UserID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     roleName,
	}, nil
}
