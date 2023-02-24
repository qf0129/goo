package crud

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateRouter[T GormModel](group *gin.RouterGroup) {
	modelName := GetModelNameLower(new(T))
	modelId := modelName + "_id"

	group.POST(fmt.Sprintf("/%s", modelName), CreateOneHandler[T]())
	group.GET(fmt.Sprintf("/%s", modelName), QueryPageHandler[T]())
	group.GET(fmt.Sprintf("/%s/:%s", modelName, modelId), QueryOneHandler[T](modelId))
	group.PUT(fmt.Sprintf("/%s/:%s", modelName, modelId), UpdateOneHandler[T](modelId))
	group.DELETE(fmt.Sprintf("/%s/:%s", modelName, modelId), DeleteOneHandler[T](modelId))
}

func CreateChildRouter[P GormModel, C GormModel](group *gin.RouterGroup) {
	parentName := GetModelNameLower(new(P))
	parentId := parentName + "_id"
	childName := GetModelNameLower(new(C))
	// childId := childName + "_id"

	group.GET(fmt.Sprintf("/%s/:%s/%s", parentName, parentId, childName), QueryPageHandler[C](parentId))
	group.POST(fmt.Sprintf("/%s/:%s/%s", parentName, parentId, childName), CreateOneHandler[C](parentId))
	// group.DELETE(fmt.Sprintf("/%s/:%s/%s/:%s", parentName, parentId, childName, childId), QueryPageChildRouter[T](modelId, childId))
}
