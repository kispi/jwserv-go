package models

import (
	"time"

	"../constants"

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

type ServiceRecordSlice []*ServiceRecord

func (s ServiceRecordSlice) Len() int      { return len(s) }
func (s ServiceRecordSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ServiceRecordSlice) Less(i, j int) bool {
	return s[i].StartedAt.Format(constants.DBTimeFormat) < s[j].StartedAt.Format(constants.DBTimeFormat)
}

// TableName TableName
func (t *ServiceRecord) TableName() string {
	return "service_records"
}

func init() {
	orm.RegisterModel(new(ServiceRecord))
}
