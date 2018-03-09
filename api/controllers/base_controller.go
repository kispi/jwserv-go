package controllers

import (
	"errors"
	"strconv"
	"strings"

	"../constants"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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
	c.Data["json"] = Response{total, data}
	c.ServeJSON()
}

// Failed Failed
func (c *BaseController) Failed(err error) {
	c.Data["json"] = Response{1, err.Error()}
	c.ServeJSON()
}

// SetQuerySeterByURIParam SetQuerySeterByURIParam
func (c *BaseController) SetQuerySeterByURIParam(qs orm.QuerySeter) (orm.QuerySeter, error) {
	uriParts := strings.Split(c.Ctx.Request.RequestURI, "?")
	if len(uriParts) == 1 {
		return nil, errors.New("no url parameters given")
	} else if uriParts[1] == "" {
		return nil, errors.New("no url parameters given")
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
			return nil, errors.New("Non exist query key")
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
