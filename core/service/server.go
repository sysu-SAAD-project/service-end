package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"github.com/sysu-saad-project/service-end/core/controller"
)

var upgrader = websocket.Upgrader{}

// GetServer return web server
func GetServer() *negroni.Negroni {
	r := mux.NewRouter()
	static := "static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))
	r.HandleFunc("/art", controller.ShowActivitiesListHandler).Methods("GET")
	r.HandleFunc("/art/{id:[0-9]+}", controller.ShowActivityDetailHandler).Methods("GET")

	s := negroni.Classic()
	s.UseHandler(r)
	return s
}
