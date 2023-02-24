package crud

import (
	"gorm.io/gorm"
)

var db *gorm.DB
var conf = defaultConf

func Init(d *gorm.DB, confs ...*Config) {
	if d == nil {
		panic("[Error] Invalid gorm.DB")
	}
	db = d

	for _, c := range confs {
		if c != nil {
			conf.Update(c)
		}
	}
}
