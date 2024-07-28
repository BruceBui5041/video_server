package common

import "time"

type SQLModel struct {
	Id        int        `json:"-" gorm:"column,id;"`
	FakeId    *UID       `json:"id" gorm:"-"`
	Status    int        `json:"status" gorm:"column:status;default:1;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column,created_at;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column,updated_at;"`
}

func (model *SQLModel) GenUID(dbType int) {
	uid := NewUID(uint32(model.Id), dbType, 1)
	model.FakeId = &uid
}
