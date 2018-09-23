package services

import (
	"sort"
	"strconv"

	"../constants"
	"../models"
)

type AreaAndRecords struct {
	Area    string
	Records []*models.ServiceRecord
}
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
	dataPerPage  = 24
	areasPerPage = 5
)

func retrieveOrNil(records []*models.ServiceRecord, i int) *models.ServiceRecord {
	defer func() *models.ServiceRecord {
		recover()
		return nil
	}()
	return records[i]
}

func createPage(areasAndRecords []*AreaAndRecords) (page [][]string) {
	// Add the first line.(Areas)
	var row []string
	for _, areaAndRecords := range areasAndRecords {
		row = append(row, areaAndRecords.Area)
		row = append(row, "")
	}
	page = append(page, row)

	for i := 0; i < dataPerPage; i++ {
		var row []string
		for _, areaAndRecords := range areasAndRecords {
			// Add LeaderName
			r := retrieveOrNil(areaAndRecords.Records, i)
			if r != nil {
				row = append(row, r.LeaderName)
			} else {
				row = append(row, "")
			}
			row = append(row, "")
		}
		page = append(page, row)

		row = []string{}
		for _, areaAndRecords := range areasAndRecords {
			// Add StartedAt, EndedAt
			r := retrieveOrNil(areaAndRecords.Records, i)
			if r != nil {
				if r.StartedAt != nil {
					row = append(row, r.StartedAt.Format(constants.DBTimeFormatDateOnly))
				} else {
					row = append(row, "")
				}
				if r.EndedAt != nil {
					row = append(row, r.EndedAt.Format(constants.DBTimeFormatDateOnly))
				} else {
					row = append(row, "")
				}
			} else {
				row = append(row, "")
				row = append(row, "")
			}
		}
		page = append(page, row)
	}
	return
}

func GeneratePages(areasAndRecords []*AreaAndRecords) [][][]string {
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

func GroupByArea(serviceRecords []*models.ServiceRecord) (records []*AreaAndRecords) {
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
		groupedSlice := models.ServiceRecordSlice(grouped[string(area)])
		sort.Sort(groupedSlice)
		record := &AreaAndRecords{
			Area:    string(area),
			Records: groupedSlice,
		}
		records = append(records, record)
	}
	return
}
