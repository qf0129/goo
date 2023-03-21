package crud

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PraseFilterOptions(c *gin.Context) ([]QueryOption, error) {
	var options []QueryOption
	for k := range c.Request.URL.Query() {
		if !ArrHasStr(FIXED_OPTIONS, k) {
			options = append(options, OptionFilterBy(k, c.Query(k)))
		}
	}
	return options, nil
}

func PraseJsonFilterOptions(c *gin.Context, jsonFieldName string) ([]QueryOption, error) {
	var options []QueryOption
	for k := range c.Request.URL.Query() {
		if !ArrHasStr(FIXED_OPTIONS, k) {
			options = append(options, OptionJsonFilterBy(jsonFieldName, k, c.Query(k)))
		}
	}
	return options, nil
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
