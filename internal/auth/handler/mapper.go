package handler

import (
	"honda-leasing-api/internal/auth"
)

func toLoginResponse(result *auth.LoginResult) LoginResponse {
	return LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		Role:         result.Role,
	}
}

func toUserProfileResponse(profile *auth.UserProfile) UserProfileResponse {
	return UserProfileResponse{
		UserID:   profile.UserID,
		Email:    profile.Email,
		FullName: profile.FullName,
		Role:     profile.Role,
	}
}
