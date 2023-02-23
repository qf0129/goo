package goo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var App *gin.Engine

func InitGin() *gin.Engine {
	gin.SetMode(Config.AppMode)
	app := gin.New()
	app.Use(gin.Logger(), gin.Recovery())
	App = app
	return app
}

func RunGin() {
	listenAddr := fmt.Sprintf("%s:%d", Config.AppHost, Config.AppPort)
	svr := &http.Server{
		Handler:      App,
		Addr:         listenAddr,
		ReadTimeout:  time.Duration(Config.AppTimeout) * time.Second,
		WriteTimeout: time.Duration(Config.AppTimeout) * time.Second,
	}
	logrus.Info("Run with " + Config.AppMode + " mode ")
	logrus.Info("Server is listening " + listenAddr)
	svr.ListenAndServe()
}
