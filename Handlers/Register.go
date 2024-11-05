package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"waterbase/Auth"
	"waterbase/DocumentDB"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		RegisterPostHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("type")

	query = strings.ToLower(query)

	switch query {
	case "service":
		RegisterService(w, r)
	case "collection":
		RegisterCollection(w, r)
	case "document":
		RegisterDocument(w, r)
	default:
		http.Error(w, "No supported register function found", http.StatusBadRequest)
		return
	}

}

func RegisterService(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	if _, ok := data["servicename"].(string); !ok {
		fmt.Println("Name not spesified in json")
		return
	}

	if _, ok := data["owner"].(string); !ok {
		fmt.Println("Owner not spesified in json")
		return
	}

	var service DocumentDB.Service

	service.Name = data["servicename"].(string)
	service.Owner = data["owner"].(string)

	if service.Name == "" || service.Owner == "" {
		http.Error(w, "Both service name and owner needs to be populated", http.StatusBadRequest)
		fmt.Println(service)
		return
	}

	jsonData := make(map[string]interface{})
	var success bool

	jsonData["auth"], success = Auth.KeyDB.CreateAuthenticationKey(service.Name, 32, rand.Intn(1000000))
	if !success {
		fmt.Println("Failed to create a auth key to the service")
		return
	}

	success = DocumentDB.DocDB.CreateNewService(service)
	if !success {
		http.Error(w, "Failed to create a new service. Maybe it exists already?", http.StatusBadGateway)
		fmt.Println("Failed to create a new service")
		return
	}

	newData, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Failed to marshal return data")
		return
	}

	Auth.KeyDB.SaveDB()
	w.Header().Add("content-type", "application/json")
	w.Write(newData)
}

func RegisterCollection(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data := make(map[string]interface{})

	fmt.Println(data)

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if _, ok := data["collectionname"].(string); !ok {
		fmt.Println("Name not spesified in json")
		return
	}

	if _, ok := data["owner"].(string); !ok {
		fmt.Println("Owner not spesified in json")
		return
	}

	if _, ok := data["servicename"].(string); !ok {
		fmt.Println("Service name not spesified in json")
		return
	}

	name := data["collectionname"].(string)
	owner := data["owner"].(string)
	servicename := data["servicename"].(string)

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		fmt.Println("The user is not authenticated to make modifications")
		fmt.Println(data)
		return
	}

	if name == "" || owner == "" || servicename == "" {
		fmt.Println("All fields must be filled to create a collection")
		return
	}

	//success := DocumentDB.DocDB.GetService(servicename).CreateNewCollection(name, owner)
	success := DocumentDB.DocDB.GetService(servicename).CreateNewCollection(name, owner)
	if !success {
		fmt.Println("Could not create collection")
		http.Error(w, "", http.StatusAlreadyReported)
		return
	}

	//DocumentDB.DocDB.GetServiceA(servicename).GetCollectionA(name).SaveCollection("./Save/" + servicename)

	http.Error(w, "", http.StatusAccepted)
}

func RegisterDocument(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data := make(map[string]interface{})

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if _, ok := data["documentname"].(string); !ok {
		fmt.Println("Name not spesified in json")
		return
	}

	if _, ok := data["owner"].(string); !ok {
		fmt.Println("Owner not spesified in json")
		return
	}

	if _, ok := data["servicename"].(string); !ok {
		fmt.Println("Service name not spesified in json")
		return
	}

	if _, ok := data["collectionname"].(string); !ok {
		fmt.Println("collection name not spesified in json")
		return
	}

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		fmt.Println("The user is not authenticated to make modifications")
		fmt.Println(data)
		return
	}

	name := data["documentname"].(string)
	owner := data["owner"].(string)
	servicename := data["servicename"].(string)
	collectionname := data["collectionname"].(string)

	if name == "" || owner == "" || servicename == "" || collectionname == "" {
		fmt.Println("All fields must be filled to create a document")
		http.Error(w, "All fields must be filled to create a document", http.StatusBadRequest)
		return
	}

	success := DocumentDB.DocDB.GetService(servicename).GetCollection(collectionname).CreateNewDocument(name, owner, data["content"])
	if !success {
		fmt.Println("Could not create the document")
		http.Error(w, "ERROR: Could not create document", http.StatusBadRequest)
		return
	}

	http.Error(w, "", http.StatusAccepted)
}
