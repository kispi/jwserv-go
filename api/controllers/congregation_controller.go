package controllers

import (
	"../models"
)

// CongregationController CongregationController
type CongregationController struct {
	BaseController
}

// Get Get
func (c *CongregationController) Get() {
	congregations := []*models.Congregation{}
	qs := models.GetModelQuerySeter(new(models.Congregation), true)
	var err error
	qs, err = c.SetQuerySeterByURIParam(qs)
	if err != nil {
		c.Failed(err)
		return
	}
	qs.All(&congregations)

	c.Success(int64(len(congregations)), congregations)
}
