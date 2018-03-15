package controllers

import (
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
