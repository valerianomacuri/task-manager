package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/valerianomacuri/task-manager/common"
	"github.com/valerianomacuri/task-manager/routers"
)

//Entry point of the program
func main() {

	// Calls startup logic
	common.StartUp()
	// Get the mux router object
	router := routers.InitRoutes()
	// Create a negroni instance
	n := negroni.Classic()
	n.UseHandler(router)
	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}
	log.Println("Listening...")
	server.ListenAndServe()

}
