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

	}
}

func RetrieveGetHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	service := query.Get("service")
	collection := query.Get("collection")
	document := query.Get("document")

	if service != "" && collection != "" && document != "" {
		GetDocument(w, r)

	} else if service != "" && collection != "" && document == "" {
		fmt.Println("ye")
		GetCollection(w, r)

	} else if service != "" && collection == "" && document == "" {
		GetService(w, r)

	}
}

func GetService(w http.ResponseWriter, r *http.Request) {

	service := r.URL.Query().Get("service")

	Sfind := DocumentDB.DocDB.GetService(service)
	if Sfind == nil {
		http.Error(w, "ERROR: Could not find the service", http.StatusBadRequest)
		return
	}

	//body, err := io.ReadAll(r.Body)
	data, err := Utils.ReadFromJSON(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		fmt.Println("Utils " + err.Error())
		return
	}

	//data := make(map[string]interface{})

	//err = json.Unmarshal(body, &data)
	//if err != nil {
	//	http.Error(w, "", http.StatusBadRequest)
	//	fmt.Println("Utils " + err.Error())
	//	return
	//}

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(data)
	if !Authenticated {
		fmt.Println("User is not authenticated")
		http.Error(w, "", http.StatusUnauthorized)
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

	service := r.URL.Query().Get("service")
	collection := r.URL.Query().Get("collection")

	fmt.Println(service)
	fmt.Println(collection)

	data, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println("Utils" + err.Error())
		return
	}

	fmt.Println(data)

	data["servicename"] = service

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(data)
	if !Authenticated {
		fmt.Println("User is not authenticated")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	yes := DocumentDB.DocDB.GetService(service).GetCollection(collection)
	if yes == nil {
		http.Error(w, "ERROR: Could not find the collection", http.StatusBadRequest)
		return
	}

	fmt.Println(*yes)

	jsonData, err := json.Marshal(&yes)
	if err != nil {
		fmt.Println("fuck off col")
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonData)
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	collection := r.URL.Query().Get("collection")
	document := r.URL.Query().Get("document")

	body, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println("Utils" + err.Error())
		return
	}

	body["servicename"] = service

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(body)
	if !Authenticated {
		fmt.Println("User is not authenticated")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

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
