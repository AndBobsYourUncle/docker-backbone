package main

import (
	"fmt"
	"github.com/AndBobsYourUncle/docker-backbone/configuration"
	"github.com/AndBobsYourUncle/docker-backbone/routing"
)

func main() {
	configuration.CheckGenerateToken()

	conf := configuration.Load()
	router := routing.GetRouter(conf)

	fmt.Println("Server up and running...")
	router.Run(":443")
}
