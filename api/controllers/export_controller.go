package controllers

import (
	"errors"
	"strings"
	"time"

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

	from := c.GetQueryValue("from")
	to := c.GetQueryValue("to")
	exportType := c.GetQueryValue("exportType")
	core.Log.Debug(from, to, exportType)

	serviceRecords := []*models.ServiceRecord{}
	core.GetModelQuerySeter(nil, new(models.ServiceRecord), false).
		Filter("congregation__id", user.Congregation.ID).
		Filter("started_at__gte", from).
		Filter("started_at__lte", to).
		Limit(-1).
		All(&serviceRecords)

	sortedRecords := services.GroupByArea(serviceRecords)
	pages := services.GeneratePages(sortedRecords)
	switch strings.ToLower(exportType) {
	case "html":
		c.Success(1, pages)
		return
	case "csv":
		c.responseAsCSV(pages)
		return
	case "excel":
		c.responseAsCSV(pages)
		return
	}
	c.Error(errors.New("UNEXPECTED_EXPORT_TYPE"))
}

func (c *ExportController) responseAsCSV(pages [][][]string) {
	var csvService services.CSVService
	fileName := "Area_" + time.Now().Format("2006-01-02_15_04_05") + ".csv"
	csv, err := csvService.NewCSV(fileName)
	if err != nil {
		core.Log.Error(err)
		return
	}
	for _, page := range pages {
		csv.AddRows(page)
	}
	fileAsByte, err := csv.SaveFileAsBytes()
	if err != nil {
		core.Log.Error(err)
		return
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Description", "File Transfer")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/csv; charset=UTF-8")
	_, err = c.Ctx.ResponseWriter.Write(fileAsByte)
	if err != nil {
		c.Error(err)
		return
	}
}
