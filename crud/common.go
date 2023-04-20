package crud

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type GormModel any

type BaseModel struct {
	Id    uint      `gorm:"primaryKey;" json:"id" form:"id"`
	Ctime time.Time `gorm:"autoCreateTime;comment:'Created Time'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime;comment:'Updated Time'" json:"utime"`
}

type DeleteMarkModel struct {
	Deleted bool `gorm:"index;default:false;" json:"deleted"`
}

type BaseUidModel struct {
	Id    uint      `gorm:"primaryKey;" json:"id" form:"id"`
	Uid   string    `gorm:"type:varchar(50);uniqueIndex;not null;" json:"uid"`
	Ctime time.Time `gorm:"autoCreateTime;comment:'Created Time'" json:"ctime"`
	Utime time.Time `gorm:"autoUpdateTime;comment:'Updated Time'" json:"utime"`
}

func (m *BaseUidModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.Uid = xid.New().String()
	return
}

// 定义查询选项类型
type QueryOption func(tx *gorm.DB) *gorm.DB

// 固定的查询选项
type FixedOption struct {
	ClosePaging bool   `form:"close_paging"` // 关闭分页，默认false
	Page        int    `form:"page"`         // 页数，默认1
	PageSize    int    `form:"page_size"`    // 每页数量，默认10
	OrderBy     string `form:"order_by"`     // 排序字段名
	Descending  bool   `form:"desc"`         // 是否倒序，默认false
	Preload     string `form:"preload"`      // 预加载表名，以英文逗号分隔
}

const (
	OPTION_CLOSE_PAGING = "close_paging"
	OPTION_PAGE         = "page"
	OPTION_PAGE_SIZE    = "page_size"
	OPTION_ORDER_BY     = "order_by"
	OPTION_DESCENDING   = "desc"
	OPTION_PRELOAD      = "preload"
)

var FIXED_OPTIONS = []string{OPTION_CLOSE_PAGING, OPTION_PAGE, OPTION_PAGE_SIZE, OPTION_ORDER_BY, OPTION_DESCENDING, OPTION_PRELOAD}
