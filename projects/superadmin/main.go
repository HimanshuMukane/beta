package main

import (
	"fmt"
	"net/http"

	"core/utils"
	"superadmin/routes"
)

func main() {
	if err := utils.LoadConfig(); err != nil {
		panic(err)
	}

	cfg := utils.AppConfig
	pool, err := utils.InitDBPool(cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	router := routes.SetupRouter(pool, cfg)
	serverPort := cfg.ServerPort

	fmt.Println("http://localhost:" + serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		panic(err)
	}
}
