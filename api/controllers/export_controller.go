package controllers

import (
	"time"

	"../core"
	"../models"
	"../services"
)

// ExportController ExportController
type ExportController struct {
	BaseController
}

// Me Me
func (c *ExportController) Me() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}
	c.Success(1, user)
}

// ExportServiceRecords exports serviceRecords
func (c *ExportController) ExportServiceRecords() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	type Payload struct {
		All   bool   `json:"all"`
		Start string `json:"start"`
		End   string `json:"end"`
	}

	payload := &Payload{}
	err = c.ParseJSONBodyStruct(&payload)
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecords := []*models.ServiceRecord{}
	core.GetModelQuerySeter(nil, new(models.ServiceRecord), false).
		Filter("congregation__id", user.Congregation.ID).
		Limit(-1).
		All(&serviceRecords)

	fileName := "Area_" + time.Now().Format("2006-01-02_15:04:05") + ".csv"
	c.Ctx.ResponseWriter.Header().Set("Content-Description", "File Transfer")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/csv")
	c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	csvService := &services.CSVService{}
	fileAsByte, err := services.Export(csvService, serviceRecords, c.Ctx.ResponseWriter, fileName)
	if err != nil {
		c.Error(err)
		return
	}
	c.Ctx.ResponseWriter.Flush()

	_, err = c.Ctx.ResponseWriter.Write(fileAsByte)
	if err != nil {
		c.Error(err)
		return
	}
}
