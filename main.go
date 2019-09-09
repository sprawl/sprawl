package main

import (
	"github.com/eqlabs/sprawl/app"
)

func main() {
	app := &app.App{}
	app.InitServices()
	app.Run()
}
