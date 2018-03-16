package main

import (
	"github.com/astaxie/beego"

	"./controllers"
)

func init() {
	userController := &controllers.UserController{}
	beego.Router("/users", userController, "get:Get")
	beego.Router("/users/:id", userController, "get:GetByID")
	beego.Router("/users/:id", userController, "put:Put")

	serviceRecordController := &controllers.ServiceRecordController{}
	beego.Router("/serviceRecords", serviceRecordController, "get:Get")
	beego.Router("/serviceRecords/:id", serviceRecordController, "get:GetByID")
	beego.Router("/serviceRecords/:id", serviceRecordController, "put:Put")

	congregationController := &controllers.CongregationController{}
	beego.Router("/congregations", congregationController, "get:Get")
	beego.Router("/congregations/:id", congregationController, "get:GetByID")
	beego.Router("/congregations/:id", congregationController, "put:Put")
}
