package main

import (
	"github.com/qf0129/goo"
	"github.com/qf0129/goo/crud"
	"github.com/qf0129/goo/middleware"
)

type Product struct {
	Id    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

func main() {
	goo.LoadCommonConfig()
	goo.LoadLogger()
	goo.ConnectDB()
	goo.MigrateModels(&Product{})

	app := goo.InitGin()
	app.Use(middleware.CorsMiddleware())

	apiGroup := app.Group("/api")

	crud.Init(goo.DB)
	crud.CreateRouter[Product](apiGroup)

	goo.RunGin()
}
