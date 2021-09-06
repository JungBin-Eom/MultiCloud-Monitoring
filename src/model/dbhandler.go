package model

import (
	"database/sql"
	"fmt"

	"github.com/JungBin-Eom/OpenStack-Logger/data"
	"github.com/JungBin-Eom/OpenStack-Logger/secret"
	_ "github.com/lib/pq" // _은 이 패키지를 명시적으로 사용하겠다는 의미
)

type postgreHandler struct {
	db *sql.DB
}

func (p *postgreHandler) Close() {
	p.db.Close()
}

func (p *postgreHandler) GetLastDate(component string) string {
	rows, err := p.db.Query(`SELECT createdOn FROM cloudlog WHERE component='` + component + `' ORDER BY createdOn DESC LIMIT 1`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	rows.Next()
	var lastTime string
	rows.Scan(&lastTime)
	return lastTime
}

func (p *postgreHandler) GetLogs(component string) []*data.Log {
	logs := []*data.Log{}
	rows, err := p.db.Query(`SELECT createdOn, component, level, message FROM cloudlog WHERE component='` + component + `' ORDER BY createdOn`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var log data.Log
		rows.Scan(&log.CreatedOn, &log.Component, &log.Level, &log.Message)
		logs = append(logs, &log)
	}
	return logs
}

func (p *postgreHandler) AddLogs(logs data.MyLog) {
	// fmt.Println("log date    : ", logs.Hits.InHits[0].Source.LogDate[0])
	// fmt.Println("component   : ", logs.Hits.InHits[0].Source.Fields.LogType)
	// fmt.Println("log type    : ", logs.Hits.InHits[0].Source.LogLevel[0])
	// fmt.Println("log message : ", logs.Hits.InHits[0].Source.LogMessage[0])
	statement, err := p.db.Prepare(`INSERT INTO cloudlog (createdOn, component, level, message) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		panic(err)
	}
	var count int
	for _, s := range logs.Hits.InHits {
		if len(s.Source.LogMessage) > 0 {
			count += 1
			_, err := statement.Exec(s.Source.LogDate, s.Source.Fields.LogType, s.Source.LogLevel[0], s.Source.LogMessage[0])
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("added ", count, "rows")
}

func (p *postgreHandler) ClearLogs(component string) bool {
	statement, err := p.db.Prepare("DELETE FROM cloudlog WHERE component=$1")
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

func (p *postgreHandler) GetError(component string) int {
	rows, err := p.db.Query("SELECT COUNT(*) FROM cloudlog WHERE component='" + component + "' AND level='ERROR'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	rows.Next()
	var count int
	rows.Scan(&count)
	return count
}

func newPostgreHandler() DBHandler {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		"cloudreamer.crrywx8kuivs.ap-northeast-2.rds.amazonaws.com", 5432, "cloudreamer", secret.GetDBPassword(), "postgres",
	)

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS cloudlog (
				createdOn DATE,
				component	TEXT,
				level 			TEXT,
				message TEXT
			);`)
	statement.Exec()
	return &postgreHandler{database}
}
