package core

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/davecgh/go-spew/spew"
)

// ModelAdapter ModelAdapter
type ModelAdapter interface {
	TableName() string
}

var (
	// UseSoftDelete UseSoftDelete
	UseSoftDelete bool
)

func init() {
	UseSoftDelete = true
}

// InsertModel wrapper of NewOrm().Insert()
func InsertModel(obj interface{}) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(obj)

	if err == nil {
		st := reflect.ValueOf(obj).Elem()
		idField := st.FieldByName("ID")
		idField.SetInt(id)
	}
	return
}

// UpdateModel wrapper of NewOrm().Update()
func UpdateModel(obj interface{}, keys []string) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(obj, keys...)

	if err != nil {
		err = errors.New(fmt.Sprintf("UpdateModel error %d %v [%s]", spew.Sdump(obj), keys, err.Error()))
	}
	return
}

// DeleteModel DeleteModel
func DeleteModel(obj interface{}) (err error) {
	if !UseSoftDelete {
		err = HardDeleteModel(obj)
	} else {
		err = SoftDeleteModel(obj)
	}
	return
}

// SoftDeleteModel mark deleted_at instead of really deletes it.
func SoftDeleteModel(m interface{}) (err error) {
	field := reflect.ValueOf(m).Elem().FieldByName("DeleteAt")
	if field.IsValid() {
		now := time.Now()
		if field.Kind() == reflect.Ptr {
			field.Set(reflect.ValueOf(&now))
		} else {
			field.Set(reflect.ValueOf(now))
		}
	}
	err = UpdateModel(m, []string{"deleted_at"})
	return
}

// HardDeleteModel really deletes data.
func HardDeleteModel(m interface{}) (err error) {
	o := orm.NewOrm()
	if _, err := o.Delete(m); err != nil {
		return err
	}
	return nil
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
