package controllers

import (
	"strconv"

	"../models"
)

// ServiceRecordController ServiceRecordController
type ServiceRecordController struct {
	BaseController
}

// Get Get
func (c *ServiceRecordController) Get() {
	serviceRecords := []*models.ServiceRecord{}
	qs := models.GetModelQuerySeter(new(models.ServiceRecord), true)
	var err error
	qs, err = c.SetQuerySeterByURIParam(qs)
	if err != nil {
		c.Failed(err)
		return
	}
	qs.All(&serviceRecords)

	c.Success(int64(len(serviceRecords)), serviceRecords)
}

// Post Post
func (c *ServiceRecordController) Post() {
	serviceRecord := new(models.ServiceRecord)
	err := c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Failed(err)
		return
	}

	err = c.PostModel(serviceRecord)
	if err != nil {
		c.Failed(err)
		return
	}

	c.Success(1, "success")
}

// GetByID GetByID
func (c *ServiceRecordController) GetByID() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Failed(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = models.GetModelQuerySeter(new(models.ServiceRecord), true).Filter("id", id).One(serviceRecord)
	if err != nil {
		c.Failed(err)
		return
	}

	c.Success(1, serviceRecord)
}

// Put Put
func (c *ServiceRecordController) Put() {
	serviceRecord := new(models.ServiceRecord)
	err := c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Failed(err)
		return
	}

	err = c.PutModel(serviceRecord)
	if err != nil {
		c.Failed(err)
		return
	}

	c.Success(1, "success")
}
