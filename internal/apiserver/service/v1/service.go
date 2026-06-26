// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

//go:generate mockgen -self_package=llmops/internal/apiserver/service/v1 -destination mock_service.go -package v1 llmops/internal/apiserver/service/v1 Service

import (
	"llmops/internal/apiserver/service/v1/user"
	useridentity "llmops/internal/apiserver/service/v1/user_identity"
	"llmops/internal/apiserver/store/mysql"
)

// Service defines functions used to return resource interface.
type Service interface {
	User() user.UserSrv
	UserIdentity() useridentity.UserIdentitySrv
}

type service struct {
	store mysql.Factory
}

// NewService returns Service interface.
func NewService(store mysql.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) User() user.UserSrv {
	return user.NewUser(s.store)
}

func (s *service) UserIdentity() useridentity.UserIdentitySrv {
	return useridentity.NewUserIdentity(s.store)
}
