package controllers

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"../constants"
	"../helper"
	"../models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
)

// BaseController UserController
type BaseController struct {
	beego.Controller
}

// Response Response
type Response struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

// Success Success
func (c *BaseController) Success(total int64, data interface{}) {
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = Response{total, data}
	c.ServeJSON()
}

// Failed Failed
func (c *BaseController) Error(err error) {
	c.Ctx.ResponseWriter.WriteHeader(500)
	c.Ctx.ResponseWriter.Write([]byte(err.Error()))
}

// ParseJSONBodyStruct ParseJSONBodyStruct
func (c *BaseController) ParseJSONBodyStruct(v interface{}) error {
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		return err
	}
	return nil
}

// ParseJSONBody ParseJSONBody
func (c *BaseController) ParseJSONBody() (json *simplejson.Json, err error) {
	json, err = simplejson.NewJson(c.Ctx.Input.RequestBody)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// SetQuerySeterByURIParam SetQuerySeterByURIParam
func (c *BaseController) SetQuerySeterByURIParam(qs orm.QuerySeter) (orm.QuerySeter, error) {
	uriParts := strings.Split(c.Ctx.Request.RequestURI, "?")
	if len(uriParts) == 1 {
		return qs, errors.New("NO_URL_PARAMETERS")
	} else if uriParts[1] == "" {
		return qs, errors.New("NO_URL_PARAMETERS")
	}
	queries := strings.Split(uriParts[1], "&")
	for _, q := range queries {
		pair := strings.Split(q, "=")
		if len(pair) != 2 {
			return nil, errors.New("wrong query format (missing '=')")
		} else if pair[1] == "" {
			return nil, errors.New("wrong query format (lack of key or value)")
		}
		switch pair[0] {
		case constants.Filter:
			var err error
			qs, err = parseFilters(qs, pair[1])
			if err != nil {
				return nil, errors.New("filter: error during query parsing")
			}
		case constants.Limit:
			val, err := strconv.ParseInt(pair[1], 10, 64)
			if err != nil {
				return nil, errors.New("limit: cannot parse value as int64")
			}
			qs = qs.Limit(val)
		case constants.Offset:
			val, err := strconv.ParseInt(pair[1], 10, 64)
			if err != nil {
				return nil, errors.New("offset: cannot parse value as int64")
			}
			qs = qs.Offset(val)
		case constants.OrderBy:
			qs = qs.OrderBy(pair[1])
		case constants.GroupBy:
			qs = qs.GroupBy(pair[1])
		default:
			return nil, errors.New("NON_EXIST_QUERY_KEY")
		}
	}
	return qs, nil
}

func parseFilters(qs orm.QuerySeter, filter string) (orm.QuerySeter, error) {
	filters := strings.Split(filter, ",")
	for _, f := range filters {
		pair := strings.Split(f, ":")
		if len(pair) != 2 {
			return nil, errors.New("wrong filter format (missing '=')")
		} else if pair[1] == "" {
			return nil, errors.New("wrong query format (lack of key or value)")
		}
		qs = qs.Filter(pair[0], pair[1])
	}
	return qs, nil
}

// PutModel PutModel
func (c *BaseController) PutModel(m interface{}) (err error) {
	idStr := c.Ctx.Input.Param(":id")
	if idStr != "" && idStr != "0" {
		id, _ := strconv.ParseInt(idStr, 10, 64)

		v := reflect.New(reflect.TypeOf(m).Elem())
		idV := v.Elem().FieldByName("ID")
		idV.SetInt(id)

		vInt := v.Interface()
		err = c.ParseJSONBodyStruct(vInt)
		if err != nil {
			return err
		}

		keys := c.GetInputKeys(vInt)
		err = models.UpdateModel(vInt, keys)
	} else {
		err = errors.New("URL parameter /:id is not given")
	}
	return
}

// GetInputKeys GetInputKeys
func (c *BaseController) GetInputKeys(v interface{}) []string {
	keysCand := helper.GetInputKeys(c.Ctx.Input.RequestBody)
	keysModels := []string{}

	vType := reflect.TypeOf(v).Elem()
	for i := 0; i < vType.NumField(); i++ {
		keyModel := vType.Field(i).Tag.Get("json")
		keyModelCleanStr := strings.Replace(keyModel, ",omitempty", "", -1)
		keysModels = append(keysModels, keyModelCleanStr)
	}

	keysUpdate := []string{}
	for _, k := range keysCand {
		if helper.HasElem(keysModels, k) {
			keysUpdate = append(keysUpdate, k)
		}
	}

	return keysUpdate
}

// GetAuthUser GetAuthUser
func (c *BaseController) GetAuthUser() (*models.User, error) {
	apikey := c.Ctx.Input.Header("apikey")
	if apikey != "" {
		authToken := new(models.AuthToken)
		err := models.GetModelQuerySeter(authToken, false).Filter("auth_token", apikey).RelatedSel("User").One(authToken)
		if err != nil {
			return nil, errors.New("Invalid Apikey")
		}
		return authToken.User, nil
	}
	return nil, errors.New("Invalid Apikey")
}
