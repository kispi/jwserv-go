package controllers

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"../core"
	"../models"
	"../services"

	"github.com/tealeg/xlsx"
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
		c.responseAsExcel(pages)
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
	fileAsByte, err := saveFileAsBytes(fileName)
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

func (c *ExportController) responseAsExcel(pages [][][]string) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("result")
	if err != nil {
		c.Error(err)
		return
	}

	for _, page := range pages {
		for _, row := range page {
			newRow := sheet.AddRow()
			for _, col := range row {
				newCell := newRow.AddCell()
				newCell.Value = col
			}
		}
	}
	filename := "temp.xlsx"
	file.Save(filename)

	c.Ctx.ResponseWriter.Header().Set("Content-Description", "File Transfer")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/csv; charset=UTF-8")
	r, err := saveFileAsBytes(filename)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = c.Ctx.ResponseWriter.Write(r)
	if err != nil {
		c.Error(err)
		return
	}
}

func saveFileAsBytes(filename string) ([]byte, error) {
	fileAsByte, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("FAILED_TO_SAVE_CREATE_BYTE_STREAM")
	}

	err = os.Remove(filename)
	if err != nil {
		core.Log.Warning("FAILED_TO_REMOVE_TEMPORARY_FILE", filename)
	}
	return fileAsByte, nil
}
