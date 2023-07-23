package crud

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/goo/pkg/arrays"
	"github.com/qf0129/goo/pkg/strs"
)

type RouterOption struct {
	Methods []string
}

func (o RouterOption) HasMethod(m string) bool {
	return arrays.HasStrItem(o.Methods, m)
}

func FilterRouterOption(options []RouterOption) RouterOption {
	var option RouterOption
	if len(options) > 0 {
		option = options[0]
	}

	if len(option.Methods) == 0 {
		option.Methods = []string{"GET", "POST", "PUT", "DELETE"}
	}
	return option
}

func CreateRouter[T GormModel](group *gin.RouterGroup, options ...RouterOption) {
	modelName := strs.GetModelNameLower(new(T))
	modelId := modelName + "_id"

	option := FilterRouterOption(options)
	if option.HasMethod("GET") {
		group.GET(fmt.Sprintf("/%s", modelName), QueryManyHandler[T]())
		group.GET(fmt.Sprintf("/%s/:%s", modelName, modelId), QueryOneHandler[T](modelId))
	}
	if option.HasMethod("POST") {
		group.POST(fmt.Sprintf("/%s", modelName), CreateOneHandler[T]())
	}
	if option.HasMethod("PUT") {
		group.PUT(fmt.Sprintf("/%s/:%s", modelName, modelId), UpdateOneHandler[T](modelId))
	}
	if option.HasMethod("DELETE") {
		group.DELETE(fmt.Sprintf("/%s/:%s", modelName, modelId), DeleteOneHandler[T](modelId))
	}
}

func CreateChildRouter[P GormModel, C GormModel](group *gin.RouterGroup, options ...RouterOption) {
	parentName := strs.GetModelNameLower(new(P))
	parentId := parentName + "_id"
	childName := strs.GetModelNameLower(new(C))
	// childId := childName + "_id"

	option := FilterRouterOption(options)
	if option.HasMethod("GET") {
		group.GET(fmt.Sprintf("/%s/:%s/%s", parentName, parentId, childName), QueryManyHandler[C](parentId))
	}
	if option.HasMethod("POST") {
		group.POST(fmt.Sprintf("/%s/:%s/%s", parentName, parentId, childName), CreateOneHandler[C](parentId))
	}
	// group.DELETE(fmt.Sprintf("/%s/:%s/%s/:%s", parentName, parentId, childName, childId), QueryPageChildRouter[T](modelId, childId))
}
