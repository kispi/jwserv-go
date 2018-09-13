package services

import (
	"bufio"
	"codebrick/core"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

type CSVService struct {
	NumCol   int
	FileName string
	Writer   *csv.Writer
}

func (c *CSVService) NewCSV(fileName string, writer io.Writer) (*CSVService, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return c, errors.New("FAILED_TO_CREATE_CSV")
	}

	c.FileName = fileName

	if writer == nil {
		writer = bufio.NewWriter(file)
	}
	wr := csv.NewWriter(writer)
	c.Writer = wr
	wr.Flush()

	return c, nil
}

func (c *CSVService) AddRow(row []string) error {
	if len(row) > c.NumCol {
		c.NumCol = len(row)
	}

	err := c.Writer.Write(row)
	if err != nil {
		return errors.New("FAILED_TO_ADD_ROW")
	}
	return nil
}

func (c *CSVService) SaveFileAsBytes() ([]byte, error) {
	fileAsByte, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return nil, errors.New("FAILED_TO_SAVE_CREATE_BYTE_STREAM")
	}

	err = os.Remove(c.FileName)
	if err != nil {
		core.Log.Warning("FAILED_TO_REMOVE_TEMPORARY_FILE", c.FileName)
	}
	return fileAsByte, nil
}
