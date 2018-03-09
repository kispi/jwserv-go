package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// BaseModel BaseModel
type BaseModel struct {
	ID        int64      `orm:"column(id);auto" json:"id"`
	CreatedAt *time.Time `orm:"column(created_at);type(timestamp)" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `orm:"column(updated_at);type(timestamp);null" json:"updatedAt,omitempty"`
	DeleteAt  *time.Time `orm:"column(deleted_at);type(timestamp);null" json:"deletedAt,omitempty"`
}

// ModelAdapter ModelAdapter
type ModelAdapter interface {
	TableName() string
}

// TableName TableName
func (m *BaseModel) TableName() string {
	return "default_table"
}

// Insert wrapper of NewOrm().Insert()
func Insert(obj interface{}) {
	o := orm.NewOrm()
	o.Using("default")

	o.Insert(obj)
}

// GetModelQuerySeter GetModelQuerySeter
func GetModelQuerySeter(m interface{}, loadRelated bool) (qs orm.QuerySeter) {
	adapter := m.(ModelAdapter)
	o := orm.NewOrm()
	qs = o.QueryTable(adapter.TableName())

	qs = ApplyCommonFilter(qs)

	if loadRelated {
		qs = qs.RelatedSel()
	}
	return
}

// ApplyCommonFilter ApplyCommonFilter
func ApplyCommonFilter(qs orm.QuerySeter) orm.QuerySeter {
	qs = qs.Filter("deleted_at__isnull", true)
	return qs
}
