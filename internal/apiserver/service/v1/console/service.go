// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

//go:generate mockgen -self_package=llmops/internal/apiserver/service/v1/console -destination mock_service.go -package console llmops/internal/apiserver/service/v1/console Service

import (
	"llmops/internal/apiserver/deps"
	"llmops/internal/apiserver/service/v1/console/user"
	useraiapikey "llmops/internal/apiserver/service/v1/console/user_ai_api_key"
	useridentity "llmops/internal/apiserver/service/v1/console/user_identity"
)

// Service defines functions used to return resource interface.
type Service interface {
	User() user.UserSrv
	UserIdentity() useridentity.UserIdentitySrv
	UserAIAPIKey() useraiapikey.UserAIAPIKeySrv
}

type service struct {
	deps *deps.Dependencies
}

// NewService returns Service interface.
func NewService(depsIns *deps.Dependencies) Service {
	return &service{
		deps: depsIns,
	}
}

func (s *service) User() user.UserSrv {
	return user.NewUser(s.deps)
}

func (s *service) UserIdentity() useridentity.UserIdentitySrv {
	return useridentity.NewUserIdentity(s.deps.MySQL)
}

func (s *service) UserAIAPIKey() useraiapikey.UserAIAPIKeySrv {
	return useraiapikey.NewUserAIAPIKey(s.deps)
}
