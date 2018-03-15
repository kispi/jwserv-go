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
	users := []*models.User{}
	qs := models.GetModelQuerySeter(new(models.User), true)
	var err error
	qs, err = c.SetQuerySeterByURIParam(qs)
	if err != nil {
		c.Failed(err)
		return
	}
	qs.All(&users)

	c.Success(int64(len(users)), users)
}
