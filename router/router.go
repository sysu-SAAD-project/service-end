package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sysu-saad-project/service-end/controller"
	"github.com/urfave/negroni"
)

var upgrader = websocket.Upgrader{}

// GetServer return web server
func GetServer() *negroni.Negroni {
	r := mux.NewRouter()
	static := "static"
	// Define static service
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	// Define /act subrouter
	act := r.PathPrefix("/act").Subrouter()
	act.HandleFunc("", controller.ShowActivitiesListHandler).Methods("GET")
	act.HandleFunc("/", controller.ShowActivitiesListHandler).Methods("GET")
	act.HandleFunc("/{id:[0-9]+}", controller.ShowActivityDetailHandler).Methods("GET")

	// Use classic server and return it
	s := negroni.Classic()
	s.UseHandler(r)
	return s
}
