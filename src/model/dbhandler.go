package model

import (
	"database/sql"

	"github.com/JungBin-Eom/OpenStack-Logger/data"
	_ "github.com/mattn/go-sqlite3" // _은 이 패키지를 명시적으로 사용하겠다는 의미
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func (s *sqliteHandler) GetLogs() []*data.Log {
	logs := []*data.Log{}
	rows, err := s.db.Query("SELECT id, created_on, type, component, message FROM logs")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var log data.Log
		rows.Scan(&log.CreatedOn, &log.Type, &log.Message)
		logs = append(logs, &log)
	}
	return logs
}

// func (s *sqliteHandler) AddLogs(logs []data.Log) {
// 	statement, err := s.db.Prepare("INSERT INTO logs (created_on, type, message) VALUES (?, ?, ?, datetime('now'))")
// 	if err != nil {
// 		panic(err)
// 	}
// 	result, err := statement.Exec(sessionId, name, false)
// 	if err != nil {
// 		panic(err)
// 	}
// 	id, _ := result.LastInsertId()
// 	var todo Todo
// 	todo.ID = int(id)
// 	todo.Name = name
// 	todo.Completed = false
// 	todo.CreatedAt = time.Now()
// 	return &todo
// }

// func (s *sqliteHandler) RemoveTodo(id int) bool {
// 	statement, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
// 	if err != nil {
// 		panic(err)
// 	}
// 	result, err := statement.Exec(id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	cnt, _ := result.RowsAffected()
// 	return cnt > 0
// }

// func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
// 	statement, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
// 	if err != nil {
// 		panic(err)
// 	}
// 	result, err := statement.Exec(complete, id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	cnt, _ := result.RowsAffected()
// 	return cnt > 0
// }

func newSqliteHandler(filepath string) DBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS logs (
				id 				INTEGER PRIMARY KEY AUTOINCREMENT,
				createdOn DATETIME,
				type 			TEXT,
				component	TEXT,
				message TEXT
			);`)
	statement.Exec()
	return &sqliteHandler{database}
}
