package main

import (
	"core"
	"especa_crm/routes"
)

func main() {
	core.Run(routes.SetupRouter)
}

