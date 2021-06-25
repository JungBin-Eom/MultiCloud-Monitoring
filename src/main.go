package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/rickyEom/logger/handlers"
)

func main() {
	l := log.New(os.Stdout, "logger", log.LstdFlags)
	lh := handlers.NewLogs(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.)
}
