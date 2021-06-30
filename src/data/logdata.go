package data

import (
	"encoding/json"
	"io"
	"time"
)

type Log struct {
	CreatedOn string `json:"created_on"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

type Logs []*Log

func GetLogs() Logs {
	return logList
}

func (l *Logs) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(l)
}

var logList = []*Log{
	&Log{
		CreatedOn: time.Now().UTC().String(),
		Type:      "DEBUG",
		Message:   "This is sample log.",
	},
}
