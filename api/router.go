package main

import (
	"github.com/astaxie/beego"

	"./controllers"
)

func init() {
	beego.Router("/users", &controllers.UserController{})
	beego.Router("/serviceRecords", &controllers.ServiceRecordController{})
	beego.Router("/congregations", &controllers.CongregationController{})
	// beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	// beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")
}
