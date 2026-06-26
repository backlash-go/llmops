// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"gorm.io/gorm"
)

// UserIdentity maps to the user_identity table.
type UserIdentity struct {
	ID        uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"column:user_id;type:bigint unsigned;not null;index:idx_user_id" json:"user_id"`
	Provider  string         `gorm:"column:provider;type:varchar(32);not null;uniqueIndex:uk_provider_subject,priority:1" json:"provider"`
	Issuer    string         `gorm:"column:issuer;type:varchar(255);not null;uniqueIndex:uk_provider_subject,priority:2" json:"issuer"`
	Subject   string         `gorm:"column:subject;type:varchar(128);not null;uniqueIndex:uk_provider_subject,priority:3" json:"subject"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"`
}

// TableName returns the database table name.
func (UserIdentity) TableName() string {
	return "user_identity"
}
