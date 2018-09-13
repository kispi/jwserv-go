package main

import (
	"github.com/astaxie/beego"

	"./controllers"
)

func init() {
	authController := &controllers.AuthController{}
	userController := &controllers.UserController{}
	serviceRecordController := &controllers.ServiceRecordController{}
	congregationController := &controllers.CongregationController{}
	exportController := &controllers.ExportController{}

	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/signUpLocal", authController, "post:SignUp"),
		beego.NSRouter("/signInLocal", authController, "post:SignIn"),

		beego.NSRouter("/me", userController, "get:Me"),
		beego.NSRouter("/users", userController, "get:Get"),
		beego.NSRouter("/users/:id", userController, "get:GetByID;put:Put;delete:Delete"),

		beego.NSRouter("/serviceRecords", serviceRecordController, "get:Get;post:Post"),
		beego.NSRouter("/serviceRecords/:id", serviceRecordController, "get:GetWithDayName;put:Put;delete:Delete"),

		beego.NSRouter("/congregations", congregationController, "get:Get;post:Post"),
		beego.NSRouter("/congregations/:id", congregationController, "get:GetByID;put:Put;delete:Delete"),

		beego.NSRouter("/export/serviceRecords", exportController, "get:ExportServiceRecords"),
	)

	beego.AddNamespace(ns)
}
