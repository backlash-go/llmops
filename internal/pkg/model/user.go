// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"gorm.io/gorm"
)

// User maps to the user table.
type User struct {
	ID          uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`
	Username    string         `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_username" json:"username"`
	Email       string         `gorm:"column:email;type:varchar(255);not null;uniqueIndex:uk_email" json:"email"`
	FirstName   string         `gorm:"column:first_name;type:varchar(64);not null;default:''" json:"first_name"`
	LastName    string         `gorm:"column:last_name;type:varchar(64);not null;default:''" json:"last_name"`
	DisplayName string         `gorm:"column:display_name;type:varchar(128);not null;default:''" json:"display_name"`
	Avatar      string         `gorm:"column:avatar;type:varchar(255);not null;default:''" json:"avatar"`
	Status      uint8          `gorm:"column:status;type:tinyint unsigned;not null;default:1" json:"status"`
	LastLoginAt *time.Time     `gorm:"column:last_login_at;type:datetime" json:"last_login_at,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"`
}

// TableName returns the database table name.
func (User) TableName() string {
	return "user"
}

// UserList is the response shape for listing users.
type UserList struct {
	metav1.ListMeta `json:",inline"`
	Items           []*User `json:"items"`
}
