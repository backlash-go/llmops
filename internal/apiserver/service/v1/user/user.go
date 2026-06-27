// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package user

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	"llmops/internal/apiserver/deps"
	"llmops/internal/pkg/code"
	"llmops/internal/pkg/model"
	"llmops/internal/pkg/session"
	apiv1 "llmops/pkg/api/llmops/v1"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, r *apiv1.CreateUserRequest) error
	OauthLogin(ctx context.Context, r *apiv1.OAuthLoginRequest) (*apiv1.OAuthLoginResponse, error)
}

type userService struct {
	deps *deps.Dependencies
}

var _ UserSrv = (*userService)(nil)

// NewUser creates a user service.
func NewUser(depsIns *deps.Dependencies) UserSrv {
	return &userService{
		deps: depsIns,
	}
}

// Create creates a user.
func (u *userService) Create(ctx context.Context, r *apiv1.CreateUserRequest) error {
	var user model.User

	_ = copier.Copy(&user, r)

	if err := u.deps.MySQL.User().Create(ctx, &user); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'uk_(username|email)'", err.Error()); match {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// OauthLogin creates or updates the local user binding after OAuth token validation.
func (u *userService) OauthLogin(ctx context.Context, r *apiv1.OAuthLoginRequest) (*apiv1.OAuthLoginResponse, error) {
	if r == nil {
		return nil, errors.WithCode(code.ErrValidation, "OAuth login request is required")
	}

	provider := strings.TrimSpace(r.Provider)
	issuer := strings.TrimSpace(r.Issuer)
	subject := strings.TrimSpace(r.Subject)
	if provider == "" || issuer == "" || subject == "" {
		return nil, errors.WithCode(code.ErrValidation, "OAuth provider, issuer and subject are required")
	}

	identity, err := u.deps.MySQL.UserIdentity().Get(ctx, &model.UserIdentity{
		Provider: provider,
		Issuer:   issuer,
		Subject:  subject,
	})

	if err == nil {
		resp, err := u.updateOAuthUser(ctx, identity, r)
		if err != nil {
			return nil, err
		}

		return u.createSession(ctx, resp)
	}

	if !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	resp, err := u.createOAuthUser(ctx, r)
	if err != nil {
		return nil, err
	}

	return u.createSession(ctx, resp)
}

func (u *userService) updateOAuthUser(
	ctx context.Context,
	identity *model.UserIdentity,
	r *apiv1.OAuthLoginRequest,
) (*apiv1.OAuthLoginResponse, error) {
	user, err := u.deps.MySQL.User().Get(ctx, &model.User{ID: identity.UserID})
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	if applyOAuthUserProfile(user, r) {
		if err := u.deps.MySQL.User().Update(ctx, user); err != nil {
			if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'uk_(username|email)'", err.Error()); match {
				return nil, errors.WithCode(code.ErrUserAlreadyExist, err.Error())
			}

			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}
	}

	return oauthLoginResponse(user, identity, false), nil
}

func (u *userService) createOAuthUser(ctx context.Context, r *apiv1.OAuthLoginRequest) (*apiv1.OAuthLoginResponse, error) {
	username := strings.TrimSpace(r.PreferredUsername)
	email := strings.TrimSpace(r.Email)
	if username == "" || email == "" {
		return nil, errors.WithCode(code.ErrValidation, "OAuth preferred_username and email are required")
	}

	now := time.Now()
	user := &model.User{
		Username:    username,
		Email:       email,
		FirstName:   strings.TrimSpace(r.FirstName),
		LastName:    strings.TrimSpace(r.LastName),
		Avatar:      strings.TrimSpace(r.Avatar),
		Status:      1,
		LastLoginAt: &now,
	}
	if err := u.deps.MySQL.User().Create(ctx, user); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'uk_(username|email)'", err.Error()); match {
			return nil, errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	identity := oauthIdentityFromRequest(r, user.ID)
	if err := u.deps.MySQL.UserIdentity().Create(ctx, identity); err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return oauthLoginResponse(user, identity, true), nil
}

func (u *userService) createSession(ctx context.Context, resp *apiv1.OAuthLoginResponse) (*apiv1.OAuthLoginResponse, error) {
	if resp == nil {
		return nil, errors.WithCode(code.ErrValidation, "OAuth login response is required")
	}
	if u.deps == nil || u.deps.Redis == nil {
		return nil, errors.WithCode(code.ErrDatabase, "redis store is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}

	sessionID, err := session.NewID()
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}

	now := time.Now()
	data := session.Data{
		UserID:     resp.UserID,
		IdentityID: resp.IdentityID,
		Username:   resp.Username,
		Email:      resp.Email,
		Provider:   resp.Provider,
		Issuer:     resp.Issuer,
		Subject:    resp.Subject,
		CreatedAt:  now.Unix(),
		ExpiresAt:  now.Add(session.DefaultTTL).Unix(),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}

	if err := u.deps.Redis.Rdb().Set(session.Key(sessionID), string(payload), session.DefaultTTL).Err(); err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	resp.SessionID = sessionID
	resp.SessionMaxAge = int(session.DefaultTTL.Seconds())

	return resp, nil
}
