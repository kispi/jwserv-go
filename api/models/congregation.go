package models

import (
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
