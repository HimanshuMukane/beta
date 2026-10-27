package main

import (
	"fmt"
	"net/http"

	"chi-mysql-server/routes"
	"chi-mysql-server/utils"
)

func main() {
	db, err := utils.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := routes.SetupRouter(db)
	serverPort := utils.AppConfig.ServerPort

	fmt.Println("http://localhost:" + serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		panic(err)
	}
}