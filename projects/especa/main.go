package main

import (
	"core"
	"especa/routes"
)

func main() {
	core.Run(routes.SetupRouter)
}
