package controllers

import (
	"errors"
	"strconv"

	"../core"
	"../models"
	"github.com/astaxie/beego/orm"
)

// CongregationController CongregationController
type CongregationController struct {
	BaseController
}

// Get Get
func (c *CongregationController) Get() {
	congregations := []*models.Congregation{}
	qs := core.GetModelQuerySeter(nil, new(models.Congregation), true)
	qs, _, _, _ = c.SetQuerySeterByURIParam(qs)
	qs.All(&congregations)

	c.Success(int64(len(congregations)), congregations)
}

// Post Post
func (c *CongregationController) Post() {
	congregation := new(models.Congregation)
	err := c.ParseJSONBodyStruct(congregation)
	if err != nil {
		c.Error(err)
		return
	}

	if core.GetModelQuerySeter(nil, new(models.Congregation), true).Filter("name", congregation.Name).Exist() {
		c.Error(errors.New("CONGREGATION_ALREADY_EXISTS"))
		return
	}

	_, err = core.InsertModel(nil, congregation)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// Delete Delete
func (c *CongregationController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	congregation := new(models.Congregation)
	err = core.GetModelQuerySeter(nil, new(models.Congregation), true).Filter("id", id).One(congregation)
	if err != nil {
		c.Error(err)
		return
	}

	err = core.DeleteModel(nil, congregation)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// GetByID GetByID
func (c *CongregationController) GetByID() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	congregation := new(models.Congregation)
	err = core.GetModelQuerySeter(nil, new(models.Congregation), true).Filter("id", id).One(congregation)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, congregation)
}

// Put Put
func (c *CongregationController) Put() {
	congregation := new(models.Congregation)
	err := c.ParseJSONBodyStruct(congregation)
	if err != nil {
		c.Error(err)
		return
	}

	o := orm.NewOrm()
	err = c.PutModel(o, congregation)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}
