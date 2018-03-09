package controllers

import (
	"github.com/astaxie/beego"
	"github.com/davecgh/go-spew/spew"

	"../models"
)

// UserController UserController
type UserController struct {
	beego.Controller
}

// Get Get
func (c *UserController) Get() {

}

// Run Run
func (c *UserController) Run() {
	records := []*models.ServiceRecord{}
	models.GetModelQuerySeter(new(models.ServiceRecord), true).Filter("congregation_id", "1").Limit(5).All(&records)
	spew.Dump(records)
}
