package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"waterbase/DocumentDB"
)

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RetrieveGetHandler(w, r)
	default:

	}
}

func RetrieveGetHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if query.Has("collection") && query.Has("service") && query.Has("document") {

		GetDocument(w, r)

	} else if query.Has("collection") && query.Has("service") {

		GetCollection(w, r)

	} else if query.Has("service") {

		GetService(w, r)

	}

	http.Error(w, "", http.StatusBadRequest)
}

func GetService(w http.ResponseWriter, r *http.Request) {

	service := r.URL.Query().Get("service")

	yes := DocumentDB.DocDB.GetService(service)
	if yes == nil {
		http.Error(w, "ERROR: Could not find the service", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(*yes)
	if err != nil {
		fmt.Println("fuck off ser")
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(data)

}

func GetCollection(w http.ResponseWriter, r *http.Request) {

	service := r.URL.Query().Get("service")
	collection := r.URL.Query().Get("collection")

	yes := DocumentDB.DocDB.GetService(service).GetCollection(collection)
	if yes == nil {
		http.Error(w, "ERROR: Could not find the collection", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(*yes)
	if err != nil {
		fmt.Println("fuck off col")
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(data)

}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	collection := r.URL.Query().Get("collection")
	document := r.URL.Query().Get("document")

	yes := DocumentDB.DocDB.GetService(service).GetCollection(collection).GetDocument(document)
	if yes == nil {
		http.Error(w, "ERROR: Could not find the document", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(*yes)
	if err != nil {
		fmt.Println("fuck off doc")
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(data)
}
