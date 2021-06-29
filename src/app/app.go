package app

import (
	"log"
	"net/http"
	"os"

	"github.com/JungBin-Eom/OpenStack-Logger/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type AppHandler struct {
	http.Handler
}

func MakeHandler() *AppHandler {
	l := log.New(os.Stdout, "openstack-logger", log.LstdFlags)
	lh := handlers.NewLogs(l)

	r := mux.NewRouter()
	neg := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))
	neg.UseHandler(r)

	a := &AppHandler{
		Handler: neg,
	}

	r.HandleFunc("/", lh.IndexHandler)
	r.HandleFunc("/sync", lh.LogSync).Methods("GET")

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	return a
}
