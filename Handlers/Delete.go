package handlers

import (
	"fmt"
	"net/http"
	"waterbase/Auth"
	"waterbase/DocumentDB"
	"waterbase/Utils"
)

func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		RemoveDeleteHandler(w, r)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

func RemoveDeleteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if !query.Has("type") {
		fmt.Println("No type is specified for deletion")
		return
	}

	switch query.Get("type") {
	case "service":
		DeleteService(w, r)
	case "collection":
		DeleteCollection(w, r)
	case "document":
		DeleteDocument(w, r)
	default:
		http.Error(w, "No", http.StatusBadRequest)
	}
}

func DeleteService(w http.ResponseWriter, r *http.Request) {

	body, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if _, ok := body["auth"].(string); !ok {
		fmt.Println("No auth key spesified")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if _, ok := body["servicename"].(string); !ok {
		fmt.Println("No service name spesified")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(body)
	if !Authenticated {
		fmt.Println("Failed to authenticate")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	success := DocumentDB.DocDB.DeleteService(body["servicename"].(string))
	if !success {
		fmt.Println("Failed to delete service")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	/*
		wd, _ := os.Getwd()

				path := filepath.Join(wd+"/Save/", body["servicename"].(string))

				fmt.Println(path)

				err = os.Remove(path + "__")
				if err != nil {
					fmt.Println(err.Error())
				}

				err = os.RemoveAll(path)
				if err != nil {
					fmt.Println(err.Error())
				}
	*/
	http.Error(w, "", http.StatusAccepted)
}

func DeleteCollection(w http.ResponseWriter, r *http.Request) {

	body, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serString := Utils.IsString(body["servicename"])
	autString := Utils.IsString(body["auth"])
	colString := Utils.IsString(body["collectionname"])

	if !serString || !autString || !colString {
		fmt.Println("Missing servicename, auth or collection")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ser := body["servicename"].(string)
	col := body["collectionname"].(string)

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(body)
	if !Authenticated {
		fmt.Println("Failed to authenticate")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	success := DocumentDB.DocDB.GetService(ser).DeleteCollection(col)
	if !success {
		fmt.Println("Failed to delete collection")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusAccepted)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {

	body, err := Utils.ReadFromJSON(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serString := Utils.IsString(body["servicename"])
	colString := Utils.IsString(body["collectionname"])
	docString := Utils.IsString(body["documentname"])
	autString := Utils.IsString(body["auth"])

	if !serString || !autString || !colString || !docString {
		fmt.Println("Missing servicename, collectionname, documentname or auth")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ser := body["servicename"].(string)
	col := body["collectionname"].(string)
	doc := body["documentname"].(string)

	Authenticated := Auth.KeyDB.CheckAuthenticationKey(body)
	if !Authenticated {
		fmt.Println("Failed to authenticate")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	success := DocumentDB.DocDB.GetService(ser).GetCollection(col).DeleteDocument(doc)
	if !success {
		fmt.Println("Failed to delete collection")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusAccepted)
}
