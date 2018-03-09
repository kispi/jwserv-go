package main

import (
	"github.com/astaxie/beego"

	"./controllers"
)

func init() {
	beego.Router("/user", &controllers.UserController{})
	// beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	// beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")
}
