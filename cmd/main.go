package main

import (
	"checkoutProject/pkg/bootstrap"
	"log"
)

func main() {
	err := bootstrap.Initialize()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := bootstrap.SetupRouter()
	r.Run("0.0.0.0:8080")
}
