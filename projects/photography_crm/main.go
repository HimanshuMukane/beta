package main

import (
	"core"
	"photography_crm/routes"
)

func main() {
	core.Run(routes.SetupRouter)
}
