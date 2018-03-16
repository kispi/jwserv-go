package controllers

import (
	"strconv"

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

// GetByID GetByID
func (c *UserController) GetByID() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Failed(err)
		return
	}

	user := new(models.User)
	err = models.GetModelQuerySeter(new(models.User), true).Filter("id", id).One(user)
	if err != nil {
		c.Failed(err)
		return
	}

	c.Success(1, user)
}

// Put Put
func (c *UserController) Put() {
	user := new(models.User)
	err := c.ParseJSONBodyStruct(user)
	if err != nil {
		c.Failed(err)
		return
	}

	err = c.PutModel(user)
	if err != nil {
		c.Failed(err)
		return
	}

	c.Success(1, "true")
}
