package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// User User
type User struct {
	BaseModel

	Nickname     string        `orm:"column(nickname);" json:"nickname,omitempty"`
	Congregation *Congregation `orm:"column(congregation_id);rel(fk);null" json:"congregation,omitempty"`
	Email        string        `orm:"column(email);" json:"email,omitempty"`
	Phone        string        `orm:"column(phone);" json:"phone,omitempty"`
	Name         string        `orm:"column(name);" json:"name,omitempty"`
	Password     string        `orm:"column(password);" json:"password,omitempty"`
	Auth         string        `orm:"column(auth);" json:"auth,omitempty"`
	LastActivity *time.Time    `orm:"column(last_activity);type(timestamp)" json:"lastActivity,omitempty"`
}

// TableName TableName
func (t *User) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(User))
}
