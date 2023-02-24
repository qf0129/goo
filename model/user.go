package model

import "github.com/qf0129/goo/crud"

type BaseUser struct {
	crud.BaseModel
	Username     string `gorm:"index;type:varchar(50)" json:"username"  form:"username"`
	PasswordHash string `gorm:"type:varchar(200)" json:"name"  form:"name"`
}
