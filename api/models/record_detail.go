package models

import (
	"../core"
	"github.com/astaxie/beego/orm"
)

// RecordDetail -
type RecordDetail struct {
	BaseModel

	Record *ServiceRecord `orm:"column(record_id);rel(fk);" json:"record"`
	Name   string         `orm:"column(name);" json:"name,omitempty"`
	Memo   string         `orm:"column(memo)" json:"memo,omitempty"`
}

// TableName TableName
func (t *RecordDetail) TableName() string {
	return "record_details"
}

func init() {
	orm.RegisterModel(new(RecordDetail))
}

// LoadRecordsDetails -
func LoadRecordsDetails(records []*ServiceRecord) {
	recordIds := []int64{}
	for _, record := range records {
		recordIds = append(recordIds, record.ID)
	}

	if len(recordIds) == 0 {
		return
	}

	recordDetails := []*RecordDetail{}
	_, err := core.GetModelQuerySeter(nil, &RecordDetail{}, false).
		Filter("record_id__in", recordIds).
		All(&recordDetails)
	if err != nil {
		core.Log.Error(err)
		return
	}

	for _, recordDetail := range recordDetails {
		for _, record := range records {
			if record.ID == recordDetail.Record.ID {
				if record.Details == nil {
					record.Details = []*RecordDetail{recordDetail}
				} else {
					record.Details = append(record.Details, recordDetail)
				}
				break
			}
		}
	}
}
