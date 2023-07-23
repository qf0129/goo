package tmpl

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id    uint      `gorm:"primaryKey;" json:"id" form:"id"`
	Ctime time.Time `gorm:"autoCreateTime;comment:'Created Time'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime;comment:'Updated Time'" json:"utime"`
}

type BaseUidModel struct {
	Id    string    `gorm:"primaryKey;type:varchar(50);" json:"id"`
	Ctime time.Time `gorm:"autoCreateTime:milli;comment:'Created Time'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime:milli;comment:'Updated Time'" json:"utime"`
}

type CommonUidModel struct {
	Id    string         `gorm:"primaryKey;type:varchar(50);" json:"id"`
	Ctime time.Time      `gorm:"autoCreateTime:milli;comment:'CreatedTime'" json:"ctime"`
	Utime time.Time      `gorm:"autoUpdateTime:milli;comment:'UpdatedTime'" json:"utime"`
	Dtime gorm.DeletedAt `gorm:"index;comment:'DeletedTime'" json:"dtime"`
}

func (m *BaseUidModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = xid.New().String()
	return
}

func (m *CommonUidModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.Id = xid.New().String()
	return
}

type User struct {
	BaseModel
	Username     string `gorm:"index;type:varchar(50)" json:"username"  form:"username"`
	PasswordHash string `gorm:"type:varchar(200)" json:"password_hash"  form:"password_hash"`
}
