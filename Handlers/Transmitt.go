package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"waterbase/Auth"
	consts "waterbase/Data"
	"waterbase/DocumentDB"
	"waterbase/Utils"
)

func TransmittHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		TransmittGetHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func TransmittGetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("type") {
	case "services":
		TransmittGetServices(w, r)
	case "collections":
		TransmittGetCollections(w, r)
	case "documents":
		TransmittGetDocuments(w, r)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}

}

func TransmittGetServices(w http.ResponseWriter, r *http.Request) {

	data := Utils.ReadHeader(r)

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	services, err := os.ReadDir(consts.DEFAULT_SAVE_LOCATION)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	var servicenames []string

	for _, h := range services {
		if !strings.Contains(h.Name(), "__") {
			servicenames = append(servicenames, h.Name())
		}
	}

	json, err := json.Marshal(servicenames)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(json)
}

func TransmittGetCollections(w http.ResponseWriter, r *http.Request) {

	data := Utils.ReadHeader(r)

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	service := DocumentDB.DocDB.GetService(data["servicename"].(string))
	if service == nil {
		fmt.Println("TRANSMITT GET: Could not find service")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var colNames []string

	files, err := os.ReadDir(consts.DEFAULT_SAVE_LOCATION + service.Name + "/")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	for _, h := range files {
		if !strings.Contains(h.Name(), "__") {
			colNames = append(colNames, h.Name())
		}
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

	data := Utils.ReadHeader(r)

	Authenticated := Auth.KeyDB.CheckForAuth(data)
	if !Authenticated {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	collection := DocumentDB.DocDB.GetService(data["servicename"].(string)).GetCollection(data["collectionname"].(string))
	if collection == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir(consts.DEFAULT_SAVE_LOCATION + data["servicename"].(string) + "/" + data["collectionname"].(string) + "/")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var docNames []string

	for _, h := range files {
		docNames = append(docNames, h.Name())
	}

	jsonData, err := json.Marshal(docNames)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonData)

}
