package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"waterbase/DocumentDB"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		RegisterPostHandler(w, r)
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

	if _, ok := data["name"].(string); !ok {
		fmt.Println("Name not spesified in json")
		return
	}

	if _, ok := data["owner"].(string); !ok {
		fmt.Println("Owner not spesified in json")
		return
	}

	var service DocumentDB.Service

	service.Name = data["name"].(string)
	service.Owner = data["owner"].(string)

	if service.Name == "" || service.Owner == "" {
		http.Error(w, "Both service name and owner needs to be populated", http.StatusBadRequest)
		fmt.Println(service)
		return
	}

	success := DocumentDB.DocDB.CreateNewService(service)
	if !success {
		http.Error(w, "Failed to create a new service. Maybe it exists already?", http.StatusBadGateway)
		fmt.Println("Failed to create a new service")
		fmt.Println(service)
		return
	}

	http.Error(w, "", http.StatusAccepted)
}

func RegisterCollection(w http.ResponseWriter, r *http.Request) {

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

	if _, ok := data["name"].(string); !ok {
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

	name := data["name"].(string)
	owner := data["owner"].(string)
	servicename := data["servicename"].(string)

	if name == "" || owner == "" || servicename == "" {
		fmt.Println("All fields must be filled to create a collection")
		return
	}

	success := DocumentDB.DocDB.GetService(servicename).CreateNewCollection(name, owner)
	if !success {
		fmt.Println("Could not create collection")
		return
	}

	http.Error(w, "", http.StatusAccepted)
}

func RegisterDocument(w http.ResponseWriter, r *http.Request) {

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

	if _, ok := data["name"].(string); !ok {
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

	name := data["name"].(string)
	owner := data["owner"].(string)
	servicename := data["servicename"].(string)
	collectionname := data["collectionname"].(string)

	if name == "" || owner == "" || servicename == "" || collectionname == "" {
		fmt.Println("All fields must be filled to create a collection")
		http.Error(w, "All fields must be filled to create a collection", http.StatusBadRequest)
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

/*
URL EXAMPLE

Register service:
http://localhost:8080/register?type=service
body JSON
{
	name: "bfvbot" string
	owner: "John" string
	auth: <key> int // Future implementation
}

Resp: HTTP Status code


Register collection:
http://localhost:8080/register?type=collection
{
	servicename: "bfvbot" string
	name: "cheaters" string
	owner: "John"
	auth: <key> int // FI
}


Register document:
http://localhost:8080/register?type=document
{
	servicename: "bfvbot" string
	collectionname: "cheaters" string
	name: "TeamKriss" string
	content: interface{}
}






*/

/*
JSON EXAMPLE

{
	name: bfvbot
	owner: John

}




*/
