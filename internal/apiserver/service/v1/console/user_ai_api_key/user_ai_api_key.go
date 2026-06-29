// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package useraiapikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	stderrors "errors"
	"strings"
	"time"

	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	"llmops/internal/apiserver/deps"
	"llmops/internal/pkg/code"
	"llmops/internal/pkg/model"
	apiv1 "llmops/pkg/api/llmops/v1"
)

const (
	apiKeyPrefixLength = 16
	apiKeyLast4Length  = 4
	apiKeyRandomBytes  = 32
	apiKeyStatusActive = 1
)

// UserAIAPIKeySrv defines functions used to handle user AI API key requests.
type UserAIAPIKeySrv interface {
	Create(ctx context.Context, r *apiv1.CreateUserAIAPIKeyRequest) (*apiv1.UserAIAPIKeyResponse, error)
	BatchDelete(ctx context.Context, r *apiv1.BatchDeleteUserAIAPIKeyRequest) (*apiv1.BatchDeleteUserAIAPIKeyResponse, error)
	Authenticate(ctx context.Context, authorization string) error
}

type userAIAPIKeyService struct {
	deps *deps.Dependencies
}

var _ UserAIAPIKeySrv = (*userAIAPIKeyService)(nil)

// NewUserAIAPIKey creates a user AI API key service.
func NewUserAIAPIKey(depsIns *deps.Dependencies) UserAIAPIKeySrv {
	return &userAIAPIKeyService{deps: depsIns}
}

// Create creates a user AI API key.
func (u *userAIAPIKeyService) Create(ctx context.Context, r *apiv1.CreateUserAIAPIKeyRequest) (*apiv1.UserAIAPIKeyResponse, error) {
	if r == nil {
		return nil, errors.WithCode(code.ErrValidation, "user AI API key request is required")
	}
	if r.UserID == 0 {
		return nil, errors.WithCode(code.ErrValidation, "user id is required")
	}

	name := strings.TrimSpace(r.Name)
	if name == "" {
		return nil, errors.WithCode(code.ErrValidation, "name is required")
	}

	if _, err := u.deps.MySQL.User().Get(ctx, &model.User{ID: r.UserID}); err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	plainKey, err := generateUserAIAPIKey()
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, err.Error())
	}

	key := &model.UserAIAPIKey{
		UserID:       r.UserID,
		Name:         name,
		APIKeyHash:   hashUserAIAPIKey(plainKey),
		APIKeyPrefix: apiKeyPrefix(plainKey),
		APIKeyLast4:  apiKeyLast4(plainKey),
		Status:       apiKeyStatusActive,
		ExpiresAt:    r.ExpiresAt,
	}

	if err := u.deps.MySQL.UserAIAPIKey().Create(ctx, key); err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	resp := userAIAPIKeyResponseFromModel(key)
	resp.APIKey = plainKey

	return resp, nil
}

// BatchDelete deletes user AI API keys.
func (u *userAIAPIKeyService) BatchDelete(ctx context.Context, r *apiv1.BatchDeleteUserAIAPIKeyRequest) (*apiv1.BatchDeleteUserAIAPIKeyResponse, error) {
	if r == nil {
		return nil, errors.WithCode(code.ErrValidation, "user AI API key delete request is required")
	}
	if r.UserID == 0 {
		return nil, errors.WithCode(code.ErrValidation, "user id is required")
	}

	ids := normalizeUserAIAPIKeyIDs(r.IDs)
	if len(ids) == 0 {
		return nil, errors.WithCode(code.ErrValidation, "ids are required")
	}

	deletedCount, err := u.deps.MySQL.UserAIAPIKey().BatchDelete(ctx, r.UserID, ids)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return &apiv1.BatchDeleteUserAIAPIKeyResponse{DeletedCount: deletedCount}, nil
}

// Authenticate validates a Higress bearer token with user AI API keys.
func (u *userAIAPIKeyService) Authenticate(ctx context.Context, authorization string) error {
	apiKey := bearerToken(authorization)
	if apiKey == "" {
		return errors.WithCode(code.ErrMissingHeader, "Authorization bearer token is required")
	}

	key, err := u.deps.MySQL.UserAIAPIKey().GetByHash(ctx, hashUserAIAPIKey(apiKey))
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrTokenInvalid, "AI API key is invalid")
		}

		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	if key.Status != apiKeyStatusActive {
		return errors.WithCode(code.ErrTokenInvalid, "AI API key is disabled")
	}
	if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
		return errors.WithCode(code.ErrExpired, "AI API key is expired")
	}

	return nil
}

func generateUserAIAPIKey() (string, error) {
	b := make([]byte, apiKeyRandomBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return "llmops_" + base64.RawURLEncoding.EncodeToString(b), nil
}

func bearerToken(authorization string) string {
	fields := strings.Fields(strings.TrimSpace(authorization))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(fields[1])
}

func hashUserAIAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))

	return hex.EncodeToString(sum[:])
}

func apiKeyPrefix(key string) string {
	if len(key) <= apiKeyPrefixLength {
		return key
	}

	return key[:apiKeyPrefixLength]
}

func apiKeyLast4(key string) string {
	if len(key) <= apiKeyLast4Length {
		return key
	}

	return key[len(key)-apiKeyLast4Length:]
}

func normalizeUserAIAPIKeyIDs(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	normalized := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}

	return normalized
}

func userAIAPIKeyResponseFromModel(key *model.UserAIAPIKey) *apiv1.UserAIAPIKeyResponse {
	if key == nil {
		return nil
	}

	return &apiv1.UserAIAPIKeyResponse{
		ID:           key.ID,
		UserID:       key.UserID,
		Name:         key.Name,
		APIKeyPrefix: key.APIKeyPrefix,
		APIKeyLast4:  key.APIKeyLast4,
		Status:       key.Status,
		ExpiresAt:    key.ExpiresAt,
		LastUsedAt:   key.LastUsedAt,
		CreatedAt:    key.CreatedAt,
		UpdatedAt:    key.UpdatedAt,
	}
}
