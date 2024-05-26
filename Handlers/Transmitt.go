package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"waterbase/Auth"
	"waterbase/DocumentDB"
	"waterbase/Utils"
)

func TransmittHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		TransmittGetHandler(w, r)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func TransmittGetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("type") {
	case "collections":
		TransmittGetCollections(w, r)
	case "documents":
		TransmittGetDocuments(w, r)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}

}

func TransmittGetCollections(w http.ResponseWriter, r *http.Request) {
	data, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	autString := Utils.IsString(data["auth"])
	serString := Utils.IsString(data["servicename"])

	if !autString || !serString {
		fmt.Println("TRANSMITT GET: Invalid data received")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(data)
	if !Authenticated {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	service := DocumentDB.DocDB.GetService(data["servicename"].(string))
	if service == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var colNames []string

	for g := range service.Collections {
		colNames = append(colNames, g)
	}

	jsonData, err := json.Marshal(colNames)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonData)
}

func TransmittGetDocuments(w http.ResponseWriter, r *http.Request) {

	data, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	autString := Utils.IsString(data["auth"])
	serString := Utils.IsString(data["servicename"])
	colString := Utils.IsString(data["collectionname"])

	if !autString || !serString || !colString {
		fmt.Println("TRANSMITT GET: Invalid data received")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(data)
	if !Authenticated {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	collection := DocumentDB.DocDB.GetService(data["servicename"].(string)).GetCollection(data["collectionname"].(string))
	if collection == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var docNames []string

	for g := range collection.Documents {
		docNames = append(docNames, g)
	}

	jsonData, err := json.Marshal(docNames)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonData)

}
