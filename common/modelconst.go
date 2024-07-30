package common

import "time"

const (
	DbTypeVideo  = 1
	DbTypeCourse = 2
	DbTypeTag    = 3
	DbTypeUser   = 4
)

const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusSuspended = "suspended"
)

const CurrentUser = "user"

type Requester interface {
	GetUserId() uint
	GetEmail() string
	GetRole() string
}

type SQLModel struct {
	Id        uint       `json:"-" gorm:"column,id;"`
	FakeId    *UID       `json:"id" gorm:"-"`
	Status    string     `json:"status" gorm:"column:status;type:ENUM('active','inactive','suspended');default:active"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column,created_at;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column,updated_at;"`
}

func (model *SQLModel) GenUID(dbType int) {
	uid := NewUID(uint32(model.Id), dbType, 1)
	model.FakeId = &uid
}
