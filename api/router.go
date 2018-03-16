package main

import (
	"github.com/astaxie/beego"

	"./controllers"
)

func init() {
	userController := &controllers.UserController{}
	beego.Router("/users", userController, "get:Get")
	beego.Router("/users/:id", userController, "get:GetByID;put:Put;delete:Delete")

	serviceRecordController := &controllers.ServiceRecordController{}
	beego.Router("/serviceRecords", serviceRecordController, "get:Get;post:Post")
	beego.Router("/serviceRecords/:id", serviceRecordController, "get:GetByID;put:Put;delete:Delete")

	congregationController := &controllers.CongregationController{}
	beego.Router("/congregations", congregationController, "get:Get;post:Post")
	beego.Router("/congregations/:id", congregationController, "get:GetByID;put:Put;delete:Delete")
}
