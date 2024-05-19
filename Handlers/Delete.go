package handlers

import "net/http"

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		DeleteGetHandler(w, r)
	}
}

func DeleteGetHandler(w http.ResponseWriter, r *http.Request) {

}
