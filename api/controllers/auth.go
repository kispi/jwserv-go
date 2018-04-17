package controllers

import (
	"errors"

	"../models"
	"golang.org/x/crypto/bcrypt"
)

// AuthController AuthController
type AuthController struct {
	BaseController
}

// SignIn SignIn
func (c *AuthController) SignIn() {
	json, err := c.ParseJSONBody()
	if err != nil {
		c.Error(err)
		return
	}
	email, _ := json.Get("email").String()
	password, _ := json.Get("password").String()
	c.signInLocal(email, password)
}

func (c *AuthController) signInLocal(email, password string) {
	cred := map[string]string{}
	cred["type"] = "local"
	cred["email"] = email
	cred["password"] = password

	c.signInWithCreds(cred)
}

// SignUp SignUp
func (c *AuthController) SignUp() {
	json, err := c.ParseJSONBody()
	if err != nil {
		c.Error(err)
		return
	}
	email, err := json.Get("email").String()
	if err != nil || email == "" {
		c.Error(errors.New("ERROR_MISSING_EMAIL"))
		return
	}
	rawPassword, err := json.Get("password").String()
	if err != nil || rawPassword == "" {
		c.Error(errors.New("ERROR_MISSING_PASSWORD"))
		return
	}
	// name, err := json.Get("name").String()
	// if err != nil || name == "" {
	// 	c.Error(errors.New("ERROR_MISSING_NAME"))
	// 	return
	// }
	// phone, err := json.Get("phone").String()
	// if err != nil || phone == "" {
	// 	c.Error(errors.New("ERROR_MISSING_PHONE"))
	// 	return
	// }

	if models.GetModelQuerySeter(new(models.User), false).Filter("email", email).Exist() {
		c.Error(errors.New("ERROR_EMAIL"))
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Error(err)
		return
	}

	user := new(models.User)
	user.Email = email
	user.Password = string(hashedBytes[:])
	user.Role = "r"
	// user.Phone = phone
	// user.Name = name

	if _, err = models.InsertModel(user); err != nil {
		c.Error(err)
		return
	}
	if count, err := models.GetModelQuerySeter(user, false).Filter("email", user.Email).Count(); err != nil {
		c.Error(errors.New("ERROR_EMAIL"))
		return
	} else if count > 1 {
		c.Error(errors.New("ERROR_EMAIL"))
		return
	}

	c.signInLocal(user.Email, rawPassword)
}

func (c *AuthController) signInWithCreds(cred map[string]string) {
	auth, err := AuthCheckLoginCallback(cred)
	if err != nil {
		c.Error(err)
	} else {
		authUser := auth
		var authToken *models.AuthToken
		if authToken, err = authUser.RenewAuthToken(); err != nil {
			c.Error(err)
		} else {
			c.Success(1, authToken)
		}
	}
}

// AuthCheckLoginCallback AuthCheckLoginCallback
func AuthCheckLoginCallback(cred map[string]string) (*models.User, error) {
	signInType := cred["type"]
	if signInType == "local" {
		user := new(models.User)
		if err := models.GetModelQuerySeter(user, false).Filter("email", cred["email"]).One(user); err == nil {
			dbPassword := []byte(user.Password)
			rawPassword := []byte(cred["password"])

			if err := bcrypt.CompareHashAndPassword(dbPassword, rawPassword); err == nil {
				return user, nil
			}
			return nil, errors.New("Invalid Password")
		}
		return nil, errors.New("Non Exist User")
	}
	return nil, errors.New("Auth Error")
}
