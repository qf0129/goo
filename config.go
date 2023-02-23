package goo

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Configuration struct {
	DbEngine   string `json:"db_engine"`
	DbHost     string `json:"db_host"`
	DbPort     uint   `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPsd      string `json:"db_psd"`
	DbDatabase string `json:"db_database"`

	AppHost    string `json:"app_host"`
	AppPort    uint   `json:"app_port"`
	AppMode    string `json:"app_mode"`
	AppTimeout uint   `json:"app_timeout"`
	LogLevel   string `json:"log_level"`
}

var Config = &Configuration{
	DbEngine:   "mysql",
	DbHost:     "127.0.0.1",
	DbPort:     3306,
	DbUser:     "root",
	DbPsd:      "root",
	DbDatabase: "database",

	AppHost:    "127.0.0.1",
	AppPort:    8080,
	AppMode:    gin.DebugMode,
	AppTimeout: 10,
	LogLevel:   "debug",
}

func LoadCommonConfig() {
	_, err := os.Stat("config.json")
	if err == nil {
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(Config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
		}
	}
}
