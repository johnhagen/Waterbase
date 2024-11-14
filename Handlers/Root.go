package handlers

import "net/http"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RootGetHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
}

func RootGetHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome to Waterbase!"))

}
