package main

import (
	"github.com/andikabahari/eoplatform/config"
	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/server/route"
)

func main() {
	app := server.NewServer(config.NewConfig())
	route.Setup(app)
	app.Run()
}
