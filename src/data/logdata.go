package data

import (
	"encoding/json"
	"io"
	"time"
)

// swagger:model
type Log struct {
	ID int `json:"id"`
	// the type for this user
	//
	// required: true
	Type        string `json:"type"`
	CreatedOn   string `json:"created_on"`
	Description string `json:"description"`
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
		ID:          1,
		Type:        "DEBUG",
		CreatedOn:   time.Now().UTC().String(),
		Description: "This is sample log.",
	},
}
