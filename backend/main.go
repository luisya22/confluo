package main

import (
	"log"

	"github.com/luisya22/confluo/backend/api"
)

func main() {
	app := api.NewApplication(api.Config{})

	err := app.Serve()
	if err != nil {
		log.Fatal(err)
	}

}
