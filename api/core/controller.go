package core

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"../constants"
	"../helpers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	simplejson "github.com/bitly/go-simplejson"
)

// Response Response
type Response struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

// Controller controller
type Controller struct {
	beego.Controller
}

// Success Success
func (c *Controller) Success(total int64, data interface{}) {
	time_now := time.Now()
	Log.Infof("%s %s %d %v", c.Ctx.Request.Method, c.Ctx.Request.RequestURI, 200, (time.Now().Sub(time_now)).String())
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = Response{total, data}
	c.ServeJSON()
}

// Failed Failed
func (c *Controller) Error(err error) {
	time_now := time.Now()
	Log.Infof("%s %s %d %v", c.Ctx.Request.Method, c.Ctx.Request.RequestURI, 500, (time.Now().Sub(time_now)).String())
	Log.Error(err)
	c.Ctx.ResponseWriter.WriteHeader(500)
	c.Ctx.ResponseWriter.Write([]byte(err.Error()))
}

// ParseJSONBodyStruct ParseJSONBodyStruct
func (c *Controller) ParseJSONBodyStruct(v interface{}) error {
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		return err
	}
	return nil
}

// ParseJSONBody ParseJSONBody
func (c *Controller) ParseJSONBody() (json *simplejson.Json, err error) {
	json, err = simplejson.NewJson(c.Ctx.Input.RequestBody)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// GetQueryValue GetQueryValue
func (c *Controller) GetQueryValue(key string) string {
	return c.Ctx.Request.URL.Query().Get(key)
}

// SetQuerySeterByURIParam SetQuerySeterByURIParam
func (c *Controller) SetQuerySeterByURIParam(qs orm.QuerySeter) (orm.QuerySeter, []string, int64, error) {
	fields := make(map[string]int64)
	uriParts := strings.Split(c.Ctx.Request.RequestURI, "?")
	if len(uriParts) == 1 {
		return qs, nil, 0, errors.New("NO_URL_PARAMETERS")
	} else if uriParts[1] == "" {
		return qs, nil, 0, errors.New("NO_URL_PARAMETERS")
	}
	queries := strings.Split(uriParts[1], "&")
	queries = helpers.MoveLimitToEnd(queries)
	var subLimit int64
	for _, q := range queries {
		pair := strings.Split(q, "=")
		if len(pair) != 2 {
			return nil, nil, 0, errors.New("wrong query format (missing '=')")
		} else if pair[1] == "" {
			return nil, nil, 0, errors.New("wrong query format (lack of key or value)")
		}
		switch pair[0] {
		case constants.Filter:
			var err error
			qs, err = parseFilters(qs, pair[1])
			if err != nil {
				return nil, nil, 0, err
			}
			fields[constants.Filter]++
		case constants.OrderBy:
			qs = qs.OrderBy(strings.Split(pair[1], ",")...)
			fields[constants.OrderBy]++
		case constants.GroupBy:
			qs = qs.GroupBy(pair[1])
			fields[constants.GroupBy]++
		case constants.Limit:
			val, err := strconv.ParseInt(pair[1], 10, 64)
			if err != nil {
				return nil, nil, 0, errors.New("limit: cannot parse value as int64")
			}
			subLimit, _ = qs.Count()
			qs = qs.Limit(val)
			fields[constants.Limit]++
		case constants.Offset:
			val, err := strconv.ParseInt(pair[1], 10, 64)
			if err != nil {
				return nil, nil, 0, errors.New("offset: cannot parse value as int64")
			}
			qs = qs.Offset(val)
			fields[constants.Offset]++
		default:
			return nil, nil, 0, errors.New("NON_EXIST_QUERY_KEY")
		}
	}

	usedFields := []string{}
	for k := range fields {
		usedFields = append(usedFields, k)
	}
	return qs, usedFields, subLimit, nil
}

// HasQueryParam returns if queryParam exists
func (c *Controller) HasQueryParam() bool {
	uriParts := strings.Split(c.Ctx.Request.RequestURI, "?")
	if len(uriParts) == 1 {
		return false
	} else if uriParts[1] == "" {
		return false
	}
	return true
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
		queryValue, _ := url.QueryUnescape(pair[1])
		qs = qs.Filter(pair[0], queryValue)
	}
	return qs, nil
}

// PutModel PutModel
func (c *Controller) PutModel(o orm.Ormer, m interface{}) (err error) {
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
		err = UpdateModel(o, vInt, keys)
	} else {
		err = errors.New("URL parameter /:id is not given")
	}
	return
}

// GetInputKeys GetInputKeys
func (c *Controller) GetInputKeys(v interface{}) []string {
	keysCand := helpers.GetInputKeys(c.Ctx.Input.RequestBody)
	keysModels := []string{}

	vType := reflect.TypeOf(v).Elem()
	for i := 0; i < vType.NumField(); i++ {
		keyModel := vType.Field(i).Tag.Get("json")
		keyModelCleanStr := strings.Replace(keyModel, ",omitempty", "", -1)
		keysModels = append(keysModels, keyModelCleanStr)
	}

	keysUpdate := []string{}
	for _, k := range keysCand {
		if helpers.HasElem(keysModels, k) {
			keysUpdate = append(keysUpdate, k)
		}
	}

	return keysUpdate
}
