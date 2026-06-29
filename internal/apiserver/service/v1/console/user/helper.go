package user

import (
	"strings"

	"llmops/internal/pkg/model"
	apiv1 "llmops/pkg/api/llmops/v1"
)

func oauthDisplayName(r *apiv1.OAuthLoginRequest) string {
	if strings.TrimSpace(r.DisplayName) != "" {
		return strings.TrimSpace(r.DisplayName)
	}

	return strings.TrimSpace(strings.TrimSpace(r.FirstName) + " " + strings.TrimSpace(r.LastName))
}

func oauthIdentityFromRequest(r *apiv1.OAuthLoginRequest, userID uint64) *model.UserIdentity {
	return &model.UserIdentity{
		UserID:   userID,
		Provider: strings.TrimSpace(r.Provider),
		Issuer:   strings.TrimSpace(r.Issuer),
		Subject:  strings.TrimSpace(r.Subject),
	}
}

func applyOAuthUserProfile(user *model.User, r *apiv1.OAuthLoginRequest) bool {
	changed := false

	if email := strings.TrimSpace(r.Email); email != "" && user.Email != email {
		user.Email = email
		changed = true
	}
	if firstName := strings.TrimSpace(r.FirstName); user.FirstName != firstName {
		user.FirstName = firstName
		changed = true
	}
	if lastName := strings.TrimSpace(r.LastName); user.LastName != lastName {
		user.LastName = lastName
		changed = true
	}

	return changed
}

func oauthLoginResponse(user *model.User, identity *model.UserIdentity, roles []string, created bool) *apiv1.OAuthLoginResponse {
	return &apiv1.OAuthLoginResponse{
		UserID:      user.ID,
		IdentityID:  identity.ID,
		Username:    user.Username,
		Email:       user.Email,
		Provider:    identity.Provider,
		Issuer:      identity.Issuer,
		Subject:     identity.Subject,
		Roles:       roles,
		Created:     created,
		LastLoginAt: user.LastLoginAt,
	}
}

func userResponseFromModel(user *model.User) *apiv1.UserResponse {
	if user == nil {
		return nil
	}

	return &apiv1.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Avatar:      user.Avatar,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
