package main

import (
	"core"
	"superadmin/routes"
)

func main() {
	core.Run(routes.SetupRouter)
}
