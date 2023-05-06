package crud

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func CreateOneHandler[T GormModel](parentIdKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(parentIdKeys) > 0 {
			var params map[string]any
			if err := c.ShouldBindJSON(&params); err != nil {
				RespFail(c, "InvalidData, "+err.Error())
				return
			}
			params[parentIdKeys[0]] = c.Param(parentIdKeys[0])

			err := CreateOne[T](params)
			if err != nil {
				RespFail(c, "CreateOneFailed, "+err.Error())
				return
			}
			RespOk(c, &params)
		} else {
			var model T
			if err := c.ShouldBindJSON(&model); err != nil {
				RespFail(c, "InvalidData, "+err.Error())
				return
			}
			err := CreateOne[T](&model)
			if err != nil {
				RespFail(c, "CreateOneFailed, "+err.Error())
				return
			}
			RespOk(c, &model)
		}
	}
}

func QueryOneHandler[T GormModel](parentIdKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ret, _ := QueryOne[T](c.Param(parentIdKey), c.Query(OPTION_PRELOAD))
		RespOk(c, ret)
	}
}

func DeleteOneHandler[T GormModel](parentIdKey string, deletedKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var deletedKey string
		if len(deletedKeys) > 0 {
			deletedKey = deletedKeys[0]
		}

		err := DeleteOne[T](c.Param(parentIdKey), deletedKey)
		if err != nil {
			if errMySQL, ok := err.(*mysql.MySQLError); ok {
				switch errMySQL.Number {
				case 1451:
					RespFail(c, "无法删除有关联数据的项")
					return
				}
			} else {
				RespFail(c, "DeleteOneFailed, "+err.Error())
				return
			}
		}
		RespOk(c, true)
	}
}

func UpdateOneHandler[T GormModel](parentIdKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var existModel T
		err := QueryOneTarget[T](c.Param(parentIdKey), &existModel)
		if err != nil {
			RespFail(c, "QueryOneTargetFailed, "+err.Error())
			return
		}

		var objMap map[string]any
		if err = c.ShouldBindJSON(&objMap); err != nil {
			RespFail(c, "InvalidData, "+err.Error())
			return
		}

		// gorm中updates结构体不支持更新空值，使用map不支持json类型
		// 因此遍历map，将子结构的map或slice转成json字符串
		for k, v := range objMap {
			valKind := reflect.ValueOf(v).Kind()
			if valKind == reflect.Map || valKind == reflect.Slice {
				bytes, err := json.Marshal(v)
				if err != nil {
					RespFail(c, "InvalidJsonValue, "+err.Error())
					return
				}
				objMap[k] = string(bytes)
			}
		}

		err = UpdateOne[T](c.Param(parentIdKey), &objMap)
		if err != nil {
			RespFail(c, "UpdateOneFailed, "+err.Error())
			return
		}

		var newModel T
		QueryOneTarget[T](c.Param(parentIdKey), &newModel)
		RespOk(c, &newModel)
	}
}

func QueryManyHandler[T GormModel](parentIdKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var fixedOptions FixedOption
		err := c.ShouldBind(&fixedOptions)
		if err != nil {
			RespFail(c, "FixedOptionError, "+err.Error())
			return
		}

		filterOptions, err := PraseFilterOptions(c)
		if err != nil {
			RespFail(c, "FilterOptionError, "+err.Error())
			return
		}

		if len(parentIdKeys) > 0 {
			option := OptionFilterBy(parentIdKeys[0], c.Param(parentIdKeys[0]))
			filterOptions = append([]QueryOption{option}, filterOptions...)
		}

		ret, err := QueryMany[T](fixedOptions, filterOptions)
		if err != nil {
			RespFail(c, "QueryFailed, "+err.Error())
			return
		}
		RespOk(c, ret)
	}
}
