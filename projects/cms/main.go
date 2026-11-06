package main

import (
	"cms/routes"
	"core"
)

func main() {
	core.Run(routes.SetupRouter)
}
