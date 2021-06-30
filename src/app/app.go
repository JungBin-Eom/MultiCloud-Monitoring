package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JungBin-Eom/OpenStack-Logger/model"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

// template-hits, outputs-hits, listcard-source, header-field, items-나머지
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

var rd *render.Render = render.New()

func (a *AppHandler) Close() {
	a.db.Close()
}

func (a *AppHandler) IndexHandler(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/index.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) GetLogs(rw http.ResponseWriter, r *http.Request) {
	logs := a.db.GetLogs()
	rd.JSON(rw, http.StatusOK, logs)
}

func (a *AppHandler) SyncLogs(rw http.ResponseWriter, r *http.Request) {
	// curl -XGET '15.164.210.67:9200/neutron-2021.06.29/_search?q=log_level:ERROR&pretty&filter_path=hits.hits._source.logmessage'
	req, err := http.NewRequest("GET", "http://15.164.210.67:9200/neutron-2021.06.30/_search?pretty&filter_path=hits.hits._source.log_date,hits.hits._source.fields,hits.hits._source.log_level,hits.hits._source.logmessage", nil)
	if err != nil {
		http.Error(rw, "Unable to get logs", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(rw, "Unable to do request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if err != nil {
		http.Error(rw, "Unable to do request", http.StatusInternalServerError)
		return
	}

	var logs MyLog
	bytes, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(bytes, &logs)
	fmt.Println("log date    : ", logs.Hits.InHits[0].Source.LogDate[0])
	fmt.Println("log type    : ", logs.Hits.InHits[0].Source.Fields.LogType)
	fmt.Println("log level   : ", logs.Hits.InHits[0].Source.LogLevel[0])
	fmt.Println("log message : ", logs.Hits.InHits[0].Source.LogMessage[0])

	rd.Text(rw, http.StatusOK, string(bytes))
}

func MakeHandler(filepath string) *AppHandler {
	r := mux.NewRouter()
	neg := negroni.Classic()
	neg.UseHandler(r)

	a := &AppHandler{
		Handler: neg,
		db:      model.NewDBHandler(filepath),
	}

	r.HandleFunc("/", a.IndexHandler)
	r.HandleFunc("/getlog", a.GetLogs).Methods("GET")
	r.HandleFunc("/sync", a.SyncLogs).Methods("GET")

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	return a
}
