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
	Cache := &CacheMem.Cache
	Cache.Init(15, 1000)
	Auth.KeyDB.Init("Keks", 200, []byte("thisis32bitlongpassphraseimusing"))
	DocumentDB.DocDB.InitDB()
	Auth.KeyDB.ReadDB()

	router := SetupRouter()

	// Start HTTP Server
	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, router))

}

func SetupRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/waterbase/transmitt", handlers.TransmittHandler)
	r.HandleFunc("/waterbase/register", handlers.RegisterHandler)
	r.HandleFunc("/waterbase/retrieve", handlers.RetrieveHandler)
	r.HandleFunc("/waterbase/remove", handlers.RemoveHandler)

	staticFileDir := http.Dir("./dashboard/")

	staticFileHandler := http.StripPrefix("/dashboard/", http.FileServer(staticFileDir))

	r.PathPrefix("/dashboard/").Handler(staticFileHandler).Methods("GET")

	wrappedRouter := corsMiddleware(r)

	return wrappedRouter
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
