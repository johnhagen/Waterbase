package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"waterbase/Auth"
	"waterbase/DocumentDB"
	"waterbase/Utils"
)

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RetrieveGetHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func RetrieveGetHandler(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Query().Get("type") {
	case "service":
		GetService(w, r)
	case "collection":
		GetCollection(w, r)
	case "document":
		GetDocument(w, r)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func GetService(w http.ResponseWriter, r *http.Request) {

	data := Utils.ReadHeader(r)

	fmt.Println(data)
	delete(data, "adminkey")

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		fmt.Println("User is not authenticated")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	service := data["servicename"].(string)

	Sfind := DocumentDB.DocDB.GetService(service)
	if Sfind == nil {
		http.Error(w, "ERROR: Could not find the service", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(*Sfind)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		fmt.Println("fuck off ser")
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonData)

}

func GetCollection(w http.ResponseWriter, r *http.Request) {

	data := Utils.ReadHeader(r)

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		fmt.Println("User is not authenticated")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	if !Utils.IsString(data["collectionname"]) {
		fmt.Println("Missing collection name")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	service := data["servicename"].(string)
	collection := data["collectionname"].(string)

	yes := DocumentDB.DocDB.GetService(service).GetCollection(collection)
	if yes == nil {
		http.Error(w, "ERROR: Could not find the collection", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(&yes)
	if err != nil {
		fmt.Println("fuck off col")
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonData)
}

func GetDocument(w http.ResponseWriter, r *http.Request) {

	serString := r.Header.Get("Servicename")
	colString := r.Header.Get("Collectionname")
	docString := r.Header.Get("Documentname")

	if serString == "" || colString == "" || docString == "" {
		fmt.Println("Missing Servicename, collectionname or documentname")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	data := Utils.ReadHeader(r)

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		fmt.Println("User is not authenticated")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	yes := DocumentDB.DocDB.GetService(data["servicename"].(string)).GetCollection(data["collectionname"].(string)).GetDocument(data["documentname"].(string))
	if yes == nil {
		http.Error(w, "ERROR: Could not find the document", http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(&yes)
	if err != nil {
		fmt.Println("fuck off doc")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(res)
}
