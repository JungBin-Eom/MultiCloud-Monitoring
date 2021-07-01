package data

import (
	"encoding/json"
	"io"
	"time"
)

type Log struct {
	ID        int    `json:"id"`
	CreatedOn string `json:"created_on"`
	Component string `json:"component"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type MyLog struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	InHits []struct {
		Source Source `json:"_source"`
	} `json:"hits"`
}

type Source struct {
	LogDate    []string `json:"log_date"`
	LogMessage []string `json:"logmessage"`
	Fields     Fields   `json:"fields"`
	LogLevel   []string `json:"log_level"`
}

type Fields struct {
	LogType string `json:"log_type"`
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
		ID:        1,
		CreatedOn: time.Now().UTC().String(),
		Component: "nova",
		Level:     "DEBUG",
		Message:   "This is sample log.",
	},
}
