package main

import (
	"github.com/astaxie/beego"

	"./controllers"
)

func init() {
	authController := &controllers.AuthController{}
	beego.Router("/signUp", authController, "post:SignUp")
	beego.Router("/signInLocal", authController, "post:SignIn")

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
