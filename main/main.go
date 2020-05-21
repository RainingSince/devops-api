package main

import (
	"cicd/api"
	"cicd/config"
	"cicd/db"
	"fmt"
)

func main() {
	// 加载配置文件 默认当前目录下的 application.yaml
	conf, err := config.LoadConfig("./application.yaml")

	if err != nil {
		fmt.Println(err)
	}

	// 加载数据库
	err = db.DBInit()

	if err != nil {
		fmt.Println(err)
	}

	// 启动 api
	err = api.ApiInit( conf.Port)

	if err != nil {
		fmt.Println(err)
	}

}
