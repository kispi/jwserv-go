package controllers

import (
	"../models"
)

// UserController UserController
type UserController struct {
	BaseController
}

// Get Get
func (c *UserController) Get() {
	records := []*models.ServiceRecord{}
	qs := models.GetModelQuerySeter(new(models.ServiceRecord), true)
	var err error
	qs, err = c.SetQuerySeterByURIParam(qs)
	if err != nil {
		c.Failed(err)
		return
	}
	qs.All(&records)

	c.Success(1, records)
}
