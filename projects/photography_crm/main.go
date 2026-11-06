package main

import (
	"chi-mysql-server/routes"
	"core"
)

func main() {
	core.Run(routes.SetupRouter)
}
