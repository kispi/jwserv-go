package controllers

import (
	"errors"
	"strconv"

	"../core"
	"../models"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

// UserController UserController
type UserController struct {
	BaseController
}

// Me Me
func (c *UserController) Me() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}
	c.Success(1, user)
}

// Get Get
func (c *UserController) Get() {
	users := []*models.User{}
	qs := core.GetModelQuerySeter(nil, new(models.User), true)
	qs, _, _, _ = c.SetQuerySeterByURIParam(qs)
	qs.All(&users)

	c.Success(int64(len(users)), users)
}

// GetByID GetByID
func (c *UserController) GetByID() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	user := new(models.User)
	err = core.GetModelQuerySeter(nil, new(models.User), true).Filter("id", id).One(user)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, user)
}

// Put Put
func (c *UserController) Put() {
	user := new(models.User)
	err := c.ParseJSONBodyStruct(user)
	if err != nil {
		c.Error(err)
		return
	}

	rawPassword := user.Password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Error(err)
		return
	}

	o := orm.NewOrm()
	core.TransBegin(o)
	user.Password = string(hashedBytes[:])
	err = core.UpdateModel(o, user, []string{"phone", "password", "role"})
	if err != nil {
		c.Error(err)
		return
	}

	if !core.GetModelQuerySeter(nil, new(models.User), true).
		Filter("congregation__id", user.Congregation.ID).
		Filter("role", "admin").
		Exist() {
		core.TransRollback(o)
		c.Error(errors.New("ADMIN_SHOULD_EXIST"))
		return
	}
	core.TransCommit(o)

	c.Success(1, "success")
}

// Delete Delete
func (c *UserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	user := new(models.User)
	err = core.GetModelQuerySeter(nil, new(models.User), true).Filter("id", id).One(user)
	if err != nil {
		c.Error(err)
		return
	}

	err = core.DeleteModel(nil, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}
