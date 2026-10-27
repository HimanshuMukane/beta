package main

import (
	"fmt"
	"net/http"

	"chi-mysql-server/routes"
	"chi-mysql-server/utils"
)

func main() {
	if err := utils.LoadConfig(); err != nil {
		panic(err)
	}

	pool, err := utils.InitDBPool()
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	router := routes.SetupRouter(pool, utils.AppConfig)
	serverPort := utils.AppConfig.ServerPort

	fmt.Println("http://localhost:" + serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		panic(err)
	}
}
