package controller

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ShowActivitiesListHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var pageNumber string
	if len(r.Form["pageNum"]) > 0 {
		pageNumber = r.Form["pageNum"][0]
	} else {
		pageNumber = "1";
	}
}

func ShowActivityDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
}