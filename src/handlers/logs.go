package handlers

import (
	"log"
	"net/http"

	"github.com/rickyEom/logger/data"
)

type MyLogs struct {
	l *log.Logger
}

func NewLogs(l *log.Logger) *MyLogs {
	return &MyLogs{l}
}

func (m *MyLogs) GetLogs(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle GET Logs")

	logs := data.GetLogs()
	err := logs.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
