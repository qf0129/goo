package crud

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func QueryPage[T GormModel](fixed FixedOption, filterOptions []QueryOption) (any, error) {
	var total int64
	var items []T

	query := db.Model(new(T))
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
			fixed.Page = conf.DefaultPageIndex
		}
		if fixed.PageSize == 0 {
			fixed.PageSize = conf.defaultPageSize
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

func QueryAll[T GormModel](fixed FixedOption, filterOptions []QueryOption) (any, error) {
	var items []T

	query := db.Model(new(T))
	for _, option := range filterOptions {
		query = option(query)
	}

	for _, field := range strings.Split(fixed.Preload, ",") {
		query = OptionPreload(field)(query)
	}
	query = OptionOrderBy(fixed.OrderBy, fixed.Descending)(query)

	ret := query.Find(&items)
	if ret.Error != nil {
		logrus.Error(ret.Error)
		return nil, ret.Error
	} else {
		return items, nil
	}
}

func QueryOne[T GormModel](modelId any, preload string) (any, error) {
	var item T
	query := db.Model(new(T)).Where(map[string]any{conf.PrimaryKey: modelId})
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
	query := db.Model(new(T)).Where(&option)
	ret := query.Take(&item)
	if ret.Error != nil {
		logrus.Error(ret.Error)
	}
	return item, ret.Error
}

func QueryOneTarget[T GormModel](modelId any, target any) error {
	return db.Model(new(T)).Where(map[string]any{conf.PrimaryKey: modelId}).Take(&target).Error
}

func CreateOne[T GormModel](obj any) error {
	return db.Create(obj).Error
}
func CreateOneWithMap[T GormModel](obj map[string]any) error {
	return db.Model(new(T)).Create(obj).Error
}

func UpdateOne[T GormModel](modelId any, params map[string]any) error {
	return db.Model(new(T)).Where(map[string]any{conf.PrimaryKey: modelId}).Updates(params).Error
}

func DeleteOne[T GormModel](modelId any) error {
	var existModel T
	err := QueryOneTarget[T](modelId, &existModel)
	if err != nil {
		return err
	}
	return db.Delete(&existModel).Error
}
