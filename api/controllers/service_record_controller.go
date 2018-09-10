package controllers

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego/orm"

	"../helpers"
	"../models"
)

// ServiceRecordController ServiceRecordController
type ServiceRecordController struct {
	BaseController
}

// Get Get
func (c *ServiceRecordController) Get() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecords := []*models.ServiceRecord{}
	qs := models.GetModelQuerySeter(new(models.ServiceRecord), true)
	qs = qs.Filter("congregation_id", user.Congregation.ID)
	qs, _, subLimit, _ := c.SetQuerySeterByURIParam(qs)
	qs.All(&serviceRecords)

	c.Success(subLimit, serviceRecords)
}

// Post Post
func (c *ServiceRecordController) Post() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	if serviceRecord.Area == "" {
		c.Error(errors.New("ERROR_MISSING_AREA"))
		return
	}

	if serviceRecord.StartedAt == nil {
		c.Error(errors.New("ERROR_MISSING_STARTED_AT"))
		return
	}

	if serviceRecord.LeaderName == "" {
		c.Error(errors.New("ERROR_MISSING_LEADER_NAME"))
		return
	}

	serviceRecord.Congregation = user.Congregation
	serviceRecord.Recorder = user

	if models.GetModelQuerySeter(new(models.ServiceRecord), false).
		Filter("congregation__id", serviceRecord.Congregation.ID).
		Filter("area", serviceRecord.Area).
		Exist() {

	}
	_, err = models.InsertModel(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// Delete Delete
func (c *ServiceRecordController) Delete() {
	user, err := c.GetAuthUser()
	if err != nil || user.Role == "public" {
		c.Error(err)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = models.GetModelQuerySeter(new(models.ServiceRecord), true).Filter("id", id).One(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	err = models.DeleteModel(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}

// GetWithDayName GetWithDayName
func (c *ServiceRecordController) GetWithDayName() {
	user, err := c.GetAuthUser()
	if err != nil {
		c.Error(err)
		return
	}

	urlParam := "id"

	arg := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		urlParam = "day"
	}

	serviceRecord := new(models.ServiceRecord)
	serviceRecords := []*models.ServiceRecord{}
	if urlParam == "id" {
		err = models.GetModelQuerySeter(new(models.ServiceRecord), true).Filter("id", id).One(serviceRecord)
		c.Success(1, serviceRecord)
		return
	} else if urlParam == "day" {
		var total int64
		days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		for _, d := range days {
			if d == arg {
				o := orm.NewOrm()
				_, err = o.Raw("SELECT id FROM service_records WHERE dayname(started_at) = ? AND congregation_id = ?", arg, user.Congregation.ID).QueryRows(&serviceRecords)
				break
			}
		}
		ids := []int64{}
		for _, r := range serviceRecords {
			ids = append(ids, r.ID)
		}
		qs := models.GetModelQuerySeter(new(models.ServiceRecord), true).Filter("id__in", ids)
		qs, fields, subTotal, _ := c.SetQuerySeterByURIParam(qs)
		total, err := qs.All(&serviceRecords)
		if err != nil || total == 0 {
			serviceRecords = []*models.ServiceRecord{}
		}

		if helpers.ContainsString(fields, "filter") {
			total = subTotal
		} else {
			total = int64(len(ids))
		}

		c.Success(total, serviceRecords)
		return
	}
	c.Error(errors.New("UNSUPPORTED_URL_PARAM"))
	return
}

// Put Put
func (c *ServiceRecordController) Put() {
	user, err := c.GetAuthUser()
	if err != nil || user.Role == "public" {
		c.Error(err)
		return
	}

	serviceRecord := new(models.ServiceRecord)
	err = c.ParseJSONBodyStruct(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	err = c.PutModel(serviceRecord)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(1, "success")
}
