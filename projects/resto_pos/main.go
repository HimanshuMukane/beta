package main

import (
	"core"
	"resto_pos/routes"
)

func main() {
	core.Run(routes.SetupRouter)
}
