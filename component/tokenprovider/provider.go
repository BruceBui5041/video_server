package tokenprovider

import (
	"errors"
	"time"
	"video_server/common"
	models "video_server/model"
)

type Provider interface {
	Generate(data TokenPayload, expiry int) (*Token, error)
	Validate(token string) (*TokenPayload, error)
}

type Token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

type TokenPayload struct {
	UserId int           `json:"user_id"`
	Roles  []models.Role `gorm:"many2many:user_roles;"`
}

var (
	ErrNotFound      = common.NewCustomError(errors.New("token not found"), "token not found", "ErrNotFound")
	ErrEncodingToken = common.NewCustomError(errors.New("error encoding token"), "error encoding token", "ErrEncodingToken")
	ErrInvalidToken  = common.NewCustomError(errors.New("invalid token"), "invalid token", "ErrInvalidTokne")
)
