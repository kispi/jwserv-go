package controllers

import (
	"errors"

	"../core"
	"../models"
)

// BaseController BaseController
type BaseController struct {
	core.Controller
}

// GetAuthUser GetAuthUser
func (c *BaseController) GetAuthUser() (*models.User, error) {
	apikey := c.Ctx.Input.Header("apikey")
	if apikey != "" {
		authToken := new(models.AuthToken)
		err := core.GetModelQuerySeter(nil, authToken, false).
			Filter("auth_token", apikey).
			RelatedSel("User").
			RelatedSel("User__Congregation").
			One(authToken)
		if err != nil {
			return nil, errors.New("INVALID_APIKEY")
		}
		return authToken.User, nil
	}
	return nil, errors.New("INVALID_APIKEY")
}
