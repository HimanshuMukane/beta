package main

import (
	"core"
	"hrms/routes"
)

func main() {
	core.Run(routes.SetupRouter)
}

