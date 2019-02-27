package services

import (
	"bufio"
	"encoding/csv"
	"errors"
	"os"

	"../core"
)

type CSVService struct {
	FileName string
	Writer   *csv.Writer
}

func (c *CSVService) NewCSV(fileName string) (*CSVService, error) {
	file, err := os.Create(fileName)
	if err != nil {
		core.Log.Warning(err)
		return c, errors.New("FAILED_TO_CREATE_CSV")
	}

	c.FileName = fileName
	wr := csv.NewWriter(bufio.NewWriter(file))
	c.Writer = wr
	wr.Flush()

	return c, nil
}

func (c *CSVService) AddRow(row []string) error {
	err := c.Writer.Write(row)
	if err != nil {
		return errors.New("FAILED_TO_ADD_ROW")
	}
	return nil
}

func (c *CSVService) AddRows(rows [][]string) error {
	err := c.Writer.WriteAll(rows)
	if err != nil {
		return errors.New("FAILED_TO_ADD_ROWS")
	}
	return nil
}
