package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"../core"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
)

// AuthToken - model
type AuthToken struct {
	BaseModel

	User      *User     `orm:"column(user_id);rel(fk);null" json:"user"`
	Token     string    `orm:"column(auth_token);null" json:"authToken"`
	LastLogin time.Time `orm:"column(last_login);type(timestamp);null" json:"lastLogin"`
}

const (
	// AuthExpireTime - after 14 days since last login, Auth Token will expire
	AuthExpireTime = float64(24 * 14)
)

// TableName - table name
func (t *AuthToken) TableName() string {
	return "auth_tokens"
}

func init() {
	orm.RegisterModel(new(AuthToken))
}

// NewAuthToken - Create new auth token
func NewAuthToken(user *User) *AuthToken {
	token := uuid.Must(uuid.NewV4()).String()
	hasher := md5.New()
	hasher.Write([]byte(token))

	return &AuthToken{
		User:      user,
		LastLogin: time.Now(),
		Token:     hex.EncodeToString(hasher.Sum(nil)),
	}
}

// ValidateAuthToken - validates auth token
func ValidateAuthToken(token string) (*AuthToken, error) {
	authToken := &AuthToken{}
	qs := core.GetModelQuerySeter(nil, new(AuthToken), true)
	err := qs.Filter("auth_token", token).Filter("user_id__deleted_at__isnull", true).One(authToken)
	if err != nil {
		return nil, nil
	}

	duration := time.Now().Sub(authToken.LastLogin)
	if duration.Hours() > AuthExpireTime {
		return nil, errors.New("AuthToken Expired. Sign In Again")
	}

	return authToken, nil
}
