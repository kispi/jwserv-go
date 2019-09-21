package models

import (
	"time"

	"../constants"

	"github.com/astaxie/beego/orm"
)

// ServiceRecord ServiceRecord
type ServiceRecord struct {
	BaseModel

	Area         string          `orm:"column(area);" json:"area"`
	StartedAt    *time.Time      `orm:"column(started_at);type(timestamp)" json:"startedAt,omitempty"`
	EndedAt      *time.Time      `orm:"column(ended_at);type(timestamp)" json:"endedAt,omitempty"`
	Congregation *Congregation   `orm:"column(congregation_id);rel(fk)" json:"congregation,omitempty"`
	LeaderName   string          `orm:"column(leader_name)" json:"leaderName,omitempty"`
	Recorder     *User           `orm:"column(recorder_id);rel(fk)" json:"recorder,omitempty"`
	Memo         string          `orm:"column(memo)" json:"memo,omitempty"`
	Details      []*RecordDetail `orm:"reverse(many);" json:"details,omitempty"`
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

// NumOfRecordsInTheSameDayForTheSameArea returns the number of records in the same day for the same area. This must be 1 always.
func (t *ServiceRecord) NumOfRecordsInTheSameDayForTheSameArea() int64 {
	type Result struct {
		Count int64 `json:"count"`
	}
	result := &Result{}
	start := t.StartedAt.Format("2006-01-02")
	o := orm.NewOrm()
	q := "SELECT COUNT(*) AS count FROM service_records WHERE " +
		"congregation_id = ? AND " +
		"area = ? AND " +
		"started_at >= ? AND started_at <= ? AND " +
		"deleted_at IS NULL"
	o.Raw(q,
		t.Congregation.ID,
		t.Area,
		start+" 00:00:00", start+" 23:59:59").QueryRow(&result)

	return result.Count
}
