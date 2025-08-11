package main

import (
	"fmt"
	"log"
	"net/http"

	app "github.com/Nathac/go-api/internals"
	"github.com/Nathac/go-api/internals/routes"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	r := routes.SetupRoutes(app)
	app.Logger.Printf("successfully created new app")
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("cannot start http connection")
	}
}
