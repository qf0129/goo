package crud

import (
	"strings"

	"github.com/qf0129/goo"
	"github.com/sirupsen/logrus"
)

func QueryMany[T GormModel](fixed FixedOption, filterOptions []QueryOption) (any, error) {
	var total int64
	var items []T

	query := goo.DB.Model(new(T))
	for _, option := range filterOptions {
		query = option(query)
	}
	query.Count(&total)

	for _, field := range strings.Split(fixed.Preload, ",") {
		query = OptionPreload(field)(query)
	}
	query = OptionOrderBy(fixed.OrderBy, fixed.Descending)(query)

	if fixed.ClosePaging {
		ret := query.Find(&items)
		if ret.Error != nil {
			logrus.Error(ret.Error)
			return nil, ret.Error
		} else {
			return items, nil
		}
	} else {
		if fixed.Page == 0 {
			fixed.Page = 1
		}
		if fixed.PageSize == 0 {
			fixed.PageSize = goo.Config.DefaultPageSize
		}
		query = OptionWithPage(fixed.Page, fixed.PageSize)(query)
		ret := query.Find(&items)
		if ret.Error != nil {
			logrus.Error(ret.Error)
			return nil, ret.Error
		} else {
			return PageBody{
				List:     &items,
				Page:     fixed.Page,
				PageSize: fixed.PageSize,
				Total:    total,
			}, nil
		}
	}
}

func QueryOne[T GormModel](modelId any, preload string) (any, error) {
	var item T
	query := goo.DB.Model(new(T)).Where(map[string]any{goo.Config.PrimaryKey: modelId})
	for _, field := range strings.Split(preload, ",") {
		query = OptionPreload(field)(query)
	}

	ret := query.Take(&item)
	if ret.Error != nil {
		logrus.Error(ret.Error)
	}
	return item, ret.Error
}

func QueryOneByMap[T GormModel](option map[string]any) (any, error) {
	var item T
	query := goo.DB.Model(new(T)).Where(option)
	ret := query.Take(&item)
	if ret.Error != nil {
		logrus.Error(ret.Error)
	}
	return item, ret.Error
}

func QueryOneTargetByMap[T GormModel](option map[string]any, target any) error {
	return goo.DB.Model(new(T)).Where(option).Take(target).Error
}

func QueryOneTarget[T GormModel](modelId any, target any) error {
	return goo.DB.Model(new(T)).Where(map[string]any{goo.Config.PrimaryKey: modelId}).Take(&target).Error
}

func CreateOne[T GormModel](obj any) error {
	return goo.DB.Model(new(T)).Create(obj).Error
}

func UpdateOne[T GormModel](modelId any, params any) error {
	return goo.DB.Model(new(T)).Where(map[string]any{goo.Config.PrimaryKey: modelId}).Updates(params).Error
}

func DeleteOne[T GormModel](modelId any) error {
	return goo.DB.Delete(new(T)).Where(map[string]any{goo.Config.PrimaryKey: modelId}).Error
}

func HasField[T GormModel](fieldKey string) bool {
	return goo.DB.Model(new(T)).Select(fieldKey).Take(new(T)).Error == nil
}
