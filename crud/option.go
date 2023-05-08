package crud

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/goo/pkg/arrays"
	"github.com/qf0129/goo/pkg/structs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PraseFilterOptions[T GormModel](c *gin.Context) ([]QueryOption, error) {
	var options []QueryOption

	var fields = structs.GetJsonFields(new(T))
	for k, v := range c.Request.URL.Query() {
		if !arrays.HasStrItem(FIXED_OPTIONS, k) && arrays.HasStrItem(fields, k) {
			kList := strings.Split(k, ":")
			if len(kList) == 2 {
				k, operater := kList[0], kList[1]
				options = append(options, GetOptionWithOperater(k, operater, v[0]))
			} else {
				options = append(options, OptionFilterBy(k, v))
			}
		}
	}
	return options, nil
}

func PraseJsonFilterOptions(c *gin.Context, jsonFieldName string) ([]QueryOption, error) {
	var options []QueryOption
	for k := range c.Request.URL.Query() {
		if !arrays.HasStrItem(FIXED_OPTIONS, k) {
			options = append(options, OptionJsonFilterBy(jsonFieldName, k, c.Query(k)))
		}
	}
	return options, nil
}

func GetOptionWithOperater(key string, operater string, val string) QueryOption {
	return func(tx *gorm.DB) *gorm.DB {
		switch operater {
		case "eq":
			return tx.Where(fmt.Sprintf("`%s` = ?", key), val)
		case "ne":
			return tx.Where(fmt.Sprintf("`%s` != ?", key), val)
		case "gt":
			return tx.Where(fmt.Sprintf("`%s` > ?", key), val)
		case "ge":
			return tx.Where(fmt.Sprintf("`%s` >= ?", key), val)
		case "lt":
			return tx.Where(fmt.Sprintf("`%s` < ?", key), val)
		case "le":
			return tx.Where(fmt.Sprintf("`%s` <= ?", key), val)
		case "in":
			return tx.Where(fmt.Sprintf("`%s` in ?", key), strings.Split(val, ","))
		case "ni":
			return tx.Where(fmt.Sprintf("`%s` not in ?", key), strings.Split(val, ","))
		case "ct":
			return tx.Where(fmt.Sprintf("`%s` like '%%%s%%'", key, val))
		case "nc":
			return tx.Where(fmt.Sprintf("`%s` not like '%%%s%%'", key, val))
		case "sw":
			return tx.Where(fmt.Sprintf("`%s` like '%s%%'", key, val))
		case "ew":
			return tx.Where(fmt.Sprintf("`%s` like '%%%s'", key, val))
		default:
			return tx.Where("? = '?'", key, val)
		}
	}
}

func OptionPreload(field string, options ...QueryOption) QueryOption {
	return func(tx *gorm.DB) *gorm.DB {
		if field == "" {
			return tx
		} else if field == "*" {
			return tx.Preload(clause.Associations)
		} else {
			return tx.Preload(cases.Title(language.Dutch).String(field), func(tx *gorm.DB) *gorm.DB {
				for _, option := range options {
					tx = option(tx)
				}
				return tx
			})
		}
	}
}

func OptionWithPage(pageIndex int, pageSize int) QueryOption {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	}
}

func OptionOrderBy(field string, descending bool) QueryOption {
	text := fmt.Sprintf("`%s`", field)
	if descending {
		text += " desc"
	}
	return func(tx *gorm.DB) *gorm.DB {
		if field == "" {
			return tx
		} else {
			return tx.Order(text)
		}
	}
}

func OptionJsonOrderBy(field string, key string, descending bool) QueryOption {
	text := fmt.Sprintf(" `%s` ->> '$.%s'", field, key)
	if descending {
		text += " desc"
	}
	return func(tx *gorm.DB) *gorm.DB {
		if key == "" {
			return tx
		} else {
			return tx.Order(text)
		}
	}
}

func OptionFilterBy(field string, value any) QueryOption {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(map[string]any{field: value})
	}
}

func OptionJsonFilterBy(field string, key string, value any) QueryOption {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(fmt.Sprintf("%s->>'$.%s'='%s'", field, key, value))
	}
}

func OptionWhere(query any, args ...any) QueryOption {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(query, args...)
	}
}
