package crud

import (
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

func DeleteOneHandler[T GormModel](parentIdKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := DeleteOne[T](c.Param(parentIdKey))
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

		var hasJson bool
		// 遍历模型字段是否包含json列
		vo := reflect.ValueOf(existModel)
		for i := 0; i < vo.NumField(); i++ {
			val := vo.Field(i).Interface()
			if reflect.TypeOf(val).String() == "datatypes.JSON" {
				hasJson = true
			}
		}

		if hasJson {
			// 若包含json列，使用结构体更新，因为可以更新json值
			var obj T
			if err = c.ShouldBindJSON(&obj); err != nil {
				RespFail(c, "InvalidData, "+err.Error())
				return
			}
			err = UpdateOne[T](c.Param(parentIdKey), &obj)
		} else {
			// 否则默认使用map更新，因为结构体不更新空值
			var objMap map[string]any
			if err = c.ShouldBindJSON(&objMap); err != nil {
				RespFail(c, "InvalidData, "+err.Error())
				return
			}
			err = UpdateOne[T](c.Param(parentIdKey), &objMap)
		}

		if err != nil {
			RespFail(c, "UpdateOneFailed, "+err.Error())
			return
		}

		var newModel T
		QueryOneTarget[T](c.Param(parentIdKey), &newModel)
		RespOk(c, &newModel)
	}
}

func QueryPageHandler[T GormModel](parentIdKeys ...string) gin.HandlerFunc {
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
