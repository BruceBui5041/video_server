package models

type Role struct {
	RoleID      uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"uniqueIndex;not null;size:50"`
	Description string
	Users       []User `gorm:"many2many:user_role;"`
}

func (Role) TableName() string {
	return "user_role"
}
