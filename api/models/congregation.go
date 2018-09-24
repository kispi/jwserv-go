package models

import (
	"../core"
	"github.com/astaxie/beego/orm"
)

// Congregation Congregation
type Congregation struct {
	BaseModel

	Name   string `orm:"column(name);" json:"name"`
	Number string `orm:"column(number);" json:"number"`
}

// TableName TableName
func (t *Congregation) TableName() string {
	return "congregations"
}

func init() {
	orm.RegisterModel(new(Congregation))
}

// AdminExists returns if admin exists for that congregation.
func (t *Congregation) AdminExists(o orm.Ormer) bool {
	return core.GetModelQuerySeter(o, new(User), false).
		Filter("role", "admin").
		Filter("congregation_id", t.ID).Exist()
}
