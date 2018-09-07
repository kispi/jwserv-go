package controllers

import (
	"errors"
	"strconv"

	"../models"
)

// ServiceRecordController ServiceRecordController
type ServiceRecordController struct {
	BaseController
}

// Get Get
func (c *ServiceRecordController) Get() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecords := []*models.ServiceRecord{}
	qs := models.GetModelQuerySeter(new(models.ServiceRecord), true)
	qs = qs.Filter("congregation_id", user.Congregation.ID)
	qs, _ = c.SetQuerySeterByURIParam(qs)
	qs.All(&serviceRecords)

	total, _ := models.GetModelQuerySeter(new(models.ServiceRecord), false).Count()
	c.Success(total, serviceRecords)
}

// Post Post
func (c *ServiceRecordController) Post() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	if serviceRecord.Area == "" {
		c.Error(errors.New("ERROR_MISSING_AREA"))
		return
	}

	if serviceRecord.StartedAt == nil {
		c.Error(errors.New("ERROR_MISSING_STARTED_AT"))
		return
	}

	if serviceRecord.LeaderName == "" {
		c.Error(errors.New("ERROR_MISSING_LEADER_NAME"))
		return
	}

	serviceRecord.Congregation = user.Congregation
	serviceRecord.Recorder = user
	_, err = models.InsertModel(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// Delete Delete
func (c *ServiceRecordController) Delete() {
	_, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = models.GetModelQuerySeter(new(models.ServiceRecord), true).Filter("id", id).One(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	err = models.DeleteModel(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// GetByID GetByID
func (c *ServiceRecordController) GetByID() {
	_, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = models.GetModelQuerySeter(new(models.ServiceRecord), true).Filter("id", id).One(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, serviceRecord)
}

// Put Put
func (c *ServiceRecordController) Put() {
	_, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	err = c.PutModel(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}
