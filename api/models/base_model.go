package models

import (
	"time"
)

// BaseModel BaseModel
type BaseModel struct {
	ID        int64      `orm:"column(id);auto" json:"id"`
	CreatedAt *time.Time `orm:"column(created_at);type(timestamp)" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `orm:"column(updated_at);type(timestamp);null" json:"updatedAt,omitempty"`
	DeleteAt  *time.Time `orm:"column(deleted_at);type(timestamp);null" json:"deletedAt,omitempty"`
}

// TableName TableName
func (m *BaseModel) TableName() string {
	return "default_table"
}
