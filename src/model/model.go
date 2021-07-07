package model

import (
	"github.com/JungBin-Eom/OpenStack-Logger/data"
)

type DBHandler interface {
	GetLastDate(string) string
	GetLogs(string) []*data.Log
	AddLogs(data.MyLog)
	ClearLogs(string) bool
	Close()
}

func NewDBHandler(filepath string) DBHandler {
	return newSqliteHandler(filepath)
}
