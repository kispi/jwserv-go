package services

import (
	"io"
	"sort"
	"strconv"

	"../constants"
	"../core"
	"../models"
)

type AreaAndRecords struct {
	Area    string
	Records []*models.ServiceRecord
}
type ExportServiceRecords struct{}
type MyString string
type MyStringSlice []MyString

// Otherwise integer ordering will not be considered.
func (s MyStringSlice) Less(i, j int) bool {
	num1, err1 := strconv.ParseInt(string(s[i]), 10, 64)
	num2, err2 := strconv.ParseInt(string(s[j]), 10, 64)
	if err1 == nil && err2 == nil {
		return num1 < num2
	}

	return s[i] < s[j]
}
func (s MyStringSlice) Len() int      { return len(s) }
func (s MyStringSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

const (
	dataPerPage  = 25
	areasPerPage = 5
)

// ExportServiceRecords exports service records
func Export(csvService *CSVService, serviceRecords []*models.ServiceRecord, writer io.Writer, fileName string) ([]byte, error) {
	s := new(ExportServiceRecords)

	csv, err := csvService.NewCSV(fileName, writer)
	if err != nil {
		return nil, err
	}
	core.Log.Debug("Export Started")
	err = s.populate(csv, serviceRecords)
	if err != nil {
		return nil, err
	}
	core.Log.Debug("Export Done")

	fileAsByte, err := csv.SaveFileAsBytes()
	if err != nil {
		return nil, err
	}

	return fileAsByte, nil
}

func (s *ExportServiceRecords) populate(csv *CSVService, serviceRecords []*models.ServiceRecord) error {
	sortedRecords := groupByArea(serviceRecords)
	pages := generatePages(csv, sortedRecords)
	for _, page := range pages {
		for _, row := range page {
			csv.AddRow(row)
		}
	}
	return nil
}

func createPage(areasAndRecords []*AreaAndRecords) (page [][]string) {
	// Add the first line.(Areas)
	var row []string
	for _, areaAndRecords := range areasAndRecords {
		row = append(row, areaAndRecords.Area)
		row = append(row, "")
	}
	page = append(page, row)

	for i := 1; i < dataPerPage; i++ {
		for _, areaAndRecords := range areasAndRecords {
			var row []string
			// Add LeaderName
			for _, r := range areaAndRecords.Records {
				row = append(row, r.LeaderName)
				row = append(row, "")
			}
			page = append(page, row)

			row = []string{}
			// Add StartedAt, EndedAt
			for _, r := range areaAndRecords.Records {
				var startedAt, endedAt string

				if r.StartedAt != nil {
					startedAt = r.StartedAt.Format(constants.DBTimeFormatDateOnly)
				}
				row = append(row, startedAt)

				if r.EndedAt != nil {
					endedAt = r.EndedAt.Format(constants.DBTimeFormatDateOnly)
				}
				row = append(row, endedAt)
			}
			page = append(page, row)
		}
	}
	return
}

func generatePages(csv *CSVService, areasAndRecords []*AreaAndRecords) [][][]string {
	var pages [][][]string
	totalPages := len(areasAndRecords)/areasPerPage + 1
	for i := 0; i < totalPages; i++ {
		if areasPerPage*(i+1) > len(areasAndRecords) {
			page := createPage(areasAndRecords[areasPerPage*i:])
			pages = append(pages, page)
			break
		}
		page := createPage(areasAndRecords[areasPerPage*i : areasPerPage*(i+1)])
		pages = append(pages, page)
	}
	return pages
}

func groupByArea(serviceRecords []*models.ServiceRecord) (records []*AreaAndRecords) {
	areasMap := make(map[string]bool)
	areas := MyStringSlice{}
	grouped := make(map[string][]*models.ServiceRecord)
	for _, r := range serviceRecords {
		areasMap[r.Area] = true
		grouped[r.Area] = append(grouped[r.Area], r)
	}

	for k := range areasMap {
		areas = append(areas, MyString(k))
	}

	sort.Sort(areas)
	for _, area := range areas {
		record := &AreaAndRecords{
			Area:    string(area),
			Records: grouped[string(area)],
		}
		records = append(records, record)
	}
	return
}
