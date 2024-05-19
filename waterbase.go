package main

import (
	"log"
	"net/http"
	"os"
	"waterbase/DocumentDB"
	handlers "waterbase/Handlers"
)

func main() {

	// Setting up port for HTTP Server
	// Get PORT env variable if present (In this case not)
	port := os.Getenv("PORT")

	// Else setup port manually
	if port == "" {
		log.Println("Preconfigured port not found: $PORT. Defaults to 8080")
		port = "8080"
	}

	DocumentDB.DocDB.InitDB()

	// Configure handler endpoints
	http.HandleFunc("/waterbase/register", handlers.RegisterHandler)
	http.HandleFunc("/waterbase/retrieve", handlers.RetrieveHandler)

	// Start HTTP Server
	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

	// Post shutdown

}
