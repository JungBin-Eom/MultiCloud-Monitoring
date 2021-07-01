package model

import (
	"database/sql"
	"fmt"

	"github.com/JungBin-Eom/OpenStack-Logger/data"
	_ "github.com/mattn/go-sqlite3" // _은 이 패키지를 명시적으로 사용하겠다는 의미
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func (s *sqliteHandler) GetLogs(component string) []*data.Log {
	logs := []*data.Log{}
	rows, err := s.db.Query("SELECT id, createdOn, component, level, message FROM openlog WHERE component=?", component)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var log data.Log
		rows.Scan(&log.ID, &log.CreatedOn, &log.Component, &log.Level, &log.Message)
		logs = append(logs, &log)
	}
	return logs
}

func (s *sqliteHandler) AddLogs(logs data.MyLog) {
	fmt.Println("log date    : ", logs.Hits.InHits[0].Source.LogDate[0])
	fmt.Println("component   : ", logs.Hits.InHits[0].Source.Fields.LogType)
	fmt.Println("log type    : ", logs.Hits.InHits[0].Source.LogLevel[0])
	fmt.Println("log message : ", logs.Hits.InHits[0].Source.LogMessage[0])
	statement, err := s.db.Prepare("INSERT INTO openlog (createdOn, component, level, message) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	for _, s := range logs.Hits.InHits {
		if len(s.Source.LogMessage) > 0 {
			result, err := statement.Exec(s.Source.LogDate[0], s.Source.Fields.LogType, s.Source.LogLevel[0], s.Source.LogMessage[0])
			if err != nil {
				panic(err)
			}
			id, _ := result.LastInsertId()
			fmt.Println("last inserted id: ", id)
		}
	}
}

func (s *sqliteHandler) ClearLogs(component string) bool {
	statement, err := s.db.Prepare("DELETE FROM openlog WHERE component=?")
	if err != nil {
		panic(err)
	}
	result, err := statement.Exec(component)
	if err != nil {
		panic(err)
	}
	count, _ := result.RowsAffected()
	return count > 0
}

func newSqliteHandler(filepath string) DBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS openlog (
				id 				INTEGER PRIMARY KEY AUTOINCREMENT,
				createdOn DATETIME,
				component	TEXT,
				level 			TEXT,
				message TEXT
			);`)
	statement.Exec()
	return &sqliteHandler{database}
}
