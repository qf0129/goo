package goo

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Configuration struct {
	DbEngine   string
	SqliteFile string

	DbHost     string
	DbPort     uint
	DbUser     string
	DbPsd      string
	DbDatabase string

	AppHost    string
	AppPort    uint
	AppMode    string
	AppTimeout uint
	LogLevel   string

	PrimaryKey      string
	DefaultPageSize int

	// 加密算法密钥
	SecretKey string

	// 令牌过期时间，单位秒
	TokenExpiredTime uint

	Custom map[string]any
}

func (c *Configuration) Set(k string, v any) {
	c.Custom[k] = v
}

func (c *Configuration) Get(k string) any {
	return c.Custom[k]
}

func (c *Configuration) Remove(k string) {
	delete(c.Custom, k)
}

var Config = &Configuration{
	DbEngine:   "sqlite",
	SqliteFile: "sqlite.db",

	DbHost:     "127.0.0.1",
	DbPort:     3306,
	DbUser:     "root",
	DbPsd:      "root",
	DbDatabase: "test",

	AppHost:    "",
	AppPort:    8080,
	AppMode:    "debug",
	AppTimeout: 60,
	LogLevel:   "debug",

	PrimaryKey:      "id",
	DefaultPageSize: 10,

	SecretKey:        "Abcd@123",
	TokenExpiredTime: 7200,

	Custom: make(map[string]any),
}

func LoadConfig() error {
	// 读取json文件
	data, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("read json file failed, err:", err)
		return err
	}

	// json数据解析到配置
	err = json.Unmarshal(data, &Config)
	if err != nil {
		fmt.Println("json unmarshal failed, err:", err)
		return err
	}

	// json数据解析到map
	var m map[string]any
	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("json unmarshal failed, err:", err)
		return err
	}

	// 加载json配置
	for k, v := range m {
		// 判断不包含的key，放到自定义配置里
		_, ok := reflect.TypeOf(*Config).FieldByName(k)
		if !ok {
			Config.Set(k, v)
		}
	}
	return nil
}
