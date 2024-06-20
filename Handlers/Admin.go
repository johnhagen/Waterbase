package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		AdminGetHandler(w, r)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}

}

func AdminGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	page, err := os.ReadFile("./Pages/index.html")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		fmt.Println("ADMIN: " + err.Error())
		return
	}

	w.Header().Add("content-type", "text/html")
	w.Write(page)
}
