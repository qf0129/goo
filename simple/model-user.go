package simple

import "github.com/qf0129/goo/crud"

type User struct {
	crud.BaseModel
	Username     string `gorm:"index;type:varchar(50)" json:"username"  form:"username"`
	PasswordHash string `gorm:"type:varchar(200)" json:"password_hash"  form:"password_hash"`
}
