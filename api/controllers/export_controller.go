package controllers

import (
	"../core"
	"../models"
	"../services"
)

// ExportController ExportController
type ExportController struct {
	BaseController
}

// ExportServiceRecords exports serviceRecords
func (c *ExportController) ExportServiceRecords() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	all := c.GetURLQueryParam("all")
	start := c.GetURLQueryParam("start")
	end := c.GetURLQueryParam("end")
	core.Log.Debug(all, start, end)

	serviceRecords := []*models.ServiceRecord{}
	core.GetModelQuerySeter(nil, new(models.ServiceRecord), false).
		Filter("congregation__id", user.Congregation.ID).
		Limit(-1).
		All(&serviceRecords)

	fileName, fileAsByte, err := services.Export(serviceRecords)
	if err != nil {
		c.Error(err)
		return
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Description", "File Transfer")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/csv")
	c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	_, err = c.Ctx.ResponseWriter.Write(fileAsByte)
	if err != nil {
		c.Error(err)
		return
	}
}
