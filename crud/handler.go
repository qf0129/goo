package crud

import (
	"time"

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

			err := CreateOneWithMap[T](params)
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
		err := QueryByID[T](c.Param(parentIdKey), &existModel)
		if err != nil {
			RespFail(c, "QueryByIDFailed, "+err.Error())
			return
		}

		var parms map[string]any
		if err := c.ShouldBindJSON(&parms); err != nil {
			RespFail(c, "InvalidData, "+err.Error())
			return
		}

		delete(parms, conf.PrimaryKey)
		err = UpdateOne[T](c.Param(parentIdKey), parms)
		if err != nil {
			RespFail(c, "UpdateOneFailed, "+err.Error())
			return
		}

		var newModel T
		QueryByID[T](c.Param(parentIdKey), &newModel)
		RespOk(c, &newModel)
	}
}

func QueryPageHandler[T GormModel](parentIdKeys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(700 * time.Microsecond)
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

		var ret any
		if fixedOptions.ClosePaging {
			ret, err = QueryAll[T](fixedOptions, filterOptions)
		} else {
			ret, err = QueryPage[T](fixedOptions, filterOptions)
		}
		if err != nil {
			RespFail(c, "QueryFailed, "+err.Error())
			return
		}
		RespOk(c, ret)
	}
}
