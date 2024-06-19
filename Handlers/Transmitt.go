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
		http.Error(w, "", http.StatusBadRequest)
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

	if r.Header.Get("Adminkey") == "" {
		http.Error(w, "", http.StatusBadRequest)
		fmt.Println("Services: No admin key filled in")
		return
	}
	data := make(map[string]interface{})
	data["adminkey"] = r.Header.Get("Adminkey")

	Authenticated := Auth.KeyDB.CheckAdminKey(data)
	if !Authenticated {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	services, err := os.ReadDir(consts.DEFAULT_SAVE)
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
	/*data, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	*/

	autString := r.Header.Get("Auth")        //Utils.IsString(data["auth"])
	serString := r.Header.Get("Servicename") //Utils.IsString(data["servicename"])

	if autString == "" || serString == "" {
		fmt.Println("TRANSMITT GET: Invalid data received")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	data := make(map[string]interface{})
	data["auth"] = autString
	data["servicename"] = serString

	/*
		Authenticated := Auth.KeyDB.CheckAuthenticationKey(data)
		if !Authenticated {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	*/

	service := DocumentDB.DocDB.GetService(data["servicename"].(string))
	if service == nil {
		fmt.Println("TRANSMITT GET: Could not find service")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var colNames []string

	files, err := os.ReadDir(consts.DEFAULT_SAVE + service.Name + "/")
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

	/*
		Authenticated := Auth.KeyDB.CheckAuthenticationKey(data)
		if !Authenticated {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	*/

	collection := DocumentDB.DocDB.GetService(data["servicename"].(string)).GetCollection(data["collectionname"].(string))
	if collection == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir(consts.DEFAULT_SAVE + data["servicename"].(string) + "/" + data["collectionname"].(string) + "/")
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
