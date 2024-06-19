package main

import (
	"log"
	"net/http"
	"os"
	"waterbase/Auth"
	CacheMem "waterbase/Cache"
	"waterbase/DocumentDB"
	handlers "waterbase/Handlers"

	"github.com/gorilla/mux"
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

	CacheMem.Cache.Init(15, 1000)
	Auth.KeyDB.Init("Keks", 100)
	DocumentDB.DocDB.InitDB()
	Auth.KeyDB.ReadDB()

	router := SetupRouter()
	// Configure handler endpoints
	//http.HandleFunc("/waterbase/transmitt", handlers.TransmittHandler)
	//http.HandleFunc("/waterbase/register", handlers.RegisterHandler)
	//http.HandleFunc("/waterbase/retrieve", handlers.RetrieveHandler)
	//http.HandleFunc("/waterbase/remove", handlers.RemoveHandler)
	//http.HandleFunc("/waterbase/admin", handlers.AdminHandler)

	// Start cache purge worker
	go CacheMem.Cache.PurgeCacheWorker(20)

	// Start HTTP Server
	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, router))

}

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/waterbase/transmitt", handlers.TransmittHandler)
	r.HandleFunc("/waterbase/register", handlers.RegisterHandler)
	r.HandleFunc("/waterbase/retrieve", handlers.RetrieveHandler)
	r.HandleFunc("/waterbase/remove", handlers.RemoveHandler)
	//r.HandleFunc("/waterbase/admin", handlers.AdminHandler)

	staticFileDir := http.Dir("./dashboard/")

	staticFileHandler := http.StripPrefix("/dashboard/", http.FileServer(staticFileDir))

	r.PathPrefix("/dashboard/").Handler(staticFileHandler).Methods("GET")
	return r
}
