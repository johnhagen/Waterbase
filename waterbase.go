package main

import (
	"log"
	"net/http"
	"os"
	"waterbase/Auth"
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

	Auth.KeyDB.Init("Keks", 100)
	DocumentDB.DocDB.InitDB()
	Auth.KeyDB.ReadDB()

	// Configure handler endpoints
	http.HandleFunc("/waterbase/transmitt", handlers.TransmittHandler)
	http.HandleFunc("/waterbase/register", handlers.RegisterHandler)
	http.HandleFunc("/waterbase/retrieve", handlers.RetrieveHandler)
	http.HandleFunc("/waterbase/remove", handlers.RemoveHandler)

	// Start HTTP Server
	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
