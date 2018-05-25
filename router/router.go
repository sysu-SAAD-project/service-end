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

	// Define generate token service
	r.HandleFunc("/token", controller.TokenHandler).Methods("GET")
	r.HandleFunc("/token/", controller.TokenHandler).Methods("GET")

	// Define /act subrouter
	act := r.PathPrefix("/act").Subrouter()
	act.HandleFunc("", controller.ShowActivitiesListHandler).Methods("GET")
	act.HandleFunc("/", controller.ShowActivitiesListHandler).Methods("GET")
	act.HandleFunc("/{id}", controller.ShowActivityDetailHandler).Methods("GET")

	// Define /users subrouter
	users := r.PathPrefix("/users").Subrouter()
	users.HandleFunc("", controller.UserLoginHandler).Methods("POST")
	users.HandleFunc("/", controller.UserLoginHandler).Methods("POST")

	// Define /actApply subrouter
	actApplys := r.PathPrefix("/actApplys").Subrouter()
	actApplys.HandleFunc("", controller.ShowActApplysListHandler).Methods("GET")
	actApplys.HandleFunc("/", controller.ShowActApplysListHandler).Methods("GET")
	actApplys.HandleFunc("/{actId}", controller.UploadActApplyHandler).Methods("POST")

	// Define /discus subrouter
	discus := r.PathPrefix("/discus").Subrouter()
	discus.HandleFunc("", controller.UploadDiscussionHandler).Methods("POST")
	discus.HandleFunc("/comments", controller.UploadCommentHandler).Methods("POST")
	discus.HandleFunc("", controller.ListDiscussionHandler).Methods("GET")
	discus.HandleFunc("/comments", controller.ListCommentsHandler).Methods("GET")

	// Use classic server and return it
	s := negroni.Classic()
	s.UseHandler(r)
	return s
}
