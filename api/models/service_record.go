package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// ServiceRecord ServiceRecord
type ServiceRecord struct {
	BaseModel

	Area         string        `orm:"column(area);" json:"area"`
	StartedAt    *time.Time    `orm:"column(started_at);type(timestamp)" json:"startedAt,omitempty"`
	EndedAt      *time.Time    `orm:"column(ended_at);type(timestamp)" json:"endedAt,omitempty"`
	Congregation *Congregation `orm:"column(congregation_id);rel(fk)" json:"congregation,omitempty"`
	LeaderName   string        `orm:"column(leader_name)" json:"leaderName,omitempty"`
	Recorder     *User         `orm:"column(recorder_id);rel(fk)" json:"recorder,omitempty"`
	Memo         string        `orm:"column(memo)" json:"memo,omitempty"`
}

// TableName TableName
func (t *ServiceRecord) TableName() string {
	return "service_records"
}

func init() {
	orm.RegisterModel(new(ServiceRecord))
}
