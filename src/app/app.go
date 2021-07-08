package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/JungBin-Eom/OpenStack-Logger/data"
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

var rd *render.Render = render.New()

func (a *AppHandler) Close() {
	a.db.Close()
}

func (a *AppHandler) IndexHandler(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/index.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) GetLogs(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	component, _ := vars["component"]
	logs := a.db.GetLogs(component)
	rd.JSON(rw, http.StatusOK, logs)
}

func (a *AppHandler) SyncLogs(rw http.ResponseWriter, r *http.Request) {
	var sync data.MyLog
	components := []string{"nova", "heat", "cinder", "neutron", "keystone", "swift"}
	for _, com := range components {
		req, err := http.NewRequest("GET", "http://15.164.210.67:9200/"+com+"/_search?pretty&filter_path=hits.hits._source.log_date,hits.hits._source.fields,hits.hits._source.log_level,hits.hits._source.logmessage", nil)
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

		var logs data.MyLog
		bytes, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(bytes, &logs)

		lastDate := a.db.GetLastDate(com)
		for _, s := range logs.Hits.InHits {
			if err != nil {
				http.Error(rw, "Unable to parse time", http.StatusInternalServerError)
			}
			if len(s.Source.LogDate) > 0 && (lastDate == "" || lastDate < s.Source.LogDate[0]) {
				sync.Hits.InHits = append(sync.Hits.InHits, s)
			}
		}
	}

	a.db.AddLogs(sync)
	rd.Text(rw, http.StatusOK, "ok")
}

func (a *AppHandler) ClearLogs(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	component, _ := vars["component"]
	req, err := http.NewRequest("DELETE", "http://15.164.210.67:9200/"+component+"/?pretty", nil)
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
	ok := a.db.ClearLogs(component)
	if ok {
		rd.Text(rw, http.StatusOK, "clear success")
	} else {
		rd.Text(rw, http.StatusOK, "clear fail")
	}
}

func (a *AppHandler) CheckLogs(rw http.ResponseWriter, r *http.Request) {
	var novaerr, heaterr, cindererr, neutronerr, keystoneerr, swifterr int
	errors := data.ComponentError{}
	novaerr = a.db.GetError("nova")
	heaterr = a.db.GetError("heat")
	cindererr = a.db.GetError("cinder")
	neutronerr = a.db.GetError("neutron")
	keystoneerr = a.db.GetError("keystone")
	swifterr = a.db.GetError("swift")

	if novaerr != 0 {
		errors.Nova = true
	}
	if heaterr != 0 {
		errors.Heat = true
	}
	if cindererr != 0 {
		errors.Cinder = true
	}
	if neutronerr != 0 {
		errors.Neutron = true
	}
	if keystoneerr != 0 {
		errors.Keystone = true
	}
	if swifterr != 0 {
		errors.Swift = true
	}

	rd.JSON(rw, http.StatusOK, errors)
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
	r.HandleFunc("/sync", a.SyncLogs).Methods("GET")
	r.HandleFunc("/{component:[a-z]+}/getlog", a.GetLogs).Methods("GET")
	r.HandleFunc("/{component:[a-z]+}/clean", a.ClearLogs).Methods("DELETE")
	r.HandleFunc("/check", a.CheckLogs).Methods("GET")

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	return a
}
