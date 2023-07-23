package main

import (
	"github.com/qf0129/goo"
	"github.com/qf0129/goo/crud"
	"github.com/qf0129/goo/tmpl"
)

type Product struct {
	Id    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

func main() {
	goo.LoadConfig()
	goo.LoadLogger()
	goo.ConnectDB()
	goo.MigrateModels(&Product{})

	app := goo.InitGin()
	app.Use(tmpl.CorsMiddleware())

	apiGroup := app.Group("/api")
	crud.CreateRouter[Product](apiGroup)

	goo.RunGin()
}
