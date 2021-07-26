package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

var rd *render.Render = render.New()

func (a *AppHandler) Close() {
	a.db.Close()
}

func (a *AppHandler) IndexHandler(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/index.html", http.StatusTemporaryRedirect)
}

// swagger:route GET /{component:[a-z]+}/getlog logs listLogs
// Returns a list of logs

// GetLogs return a list of logs
func (a *AppHandler) GetLogs(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	component, _ := vars["component"]
	logs := a.db.GetLogs(component)
	rd.JSON(rw, http.StatusOK, logs)
}

func (a *AppHandler) SyncLogs(rw http.ResponseWriter, r *http.Request) {
	var sync data.MyLog
	components := []string{"nova", "heat", "cinder", "neutron", "keystone", "agent", "management"}
	for _, com := range components {
		req, err := http.NewRequest("GET", "http://15.164.210.67:9200/"+com+"/_search?pretty&filter_path=hits.hits._source.log_date,hits.hits._source.fields,hits.hits._source.log_level,hits.hits._source.logmessage", nil)
		if err != nil {
			http.Error(rw, "Unable to get logs", http.StatusBadRequest)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(rw, "Unable to do request", http.StatusBadRequest)
			return
		}
		defer res.Body.Close()

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
		http.Error(rw, "Unable to delete elasticsearch index", http.StatusBadRequest)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(rw, "Unable to do request", http.StatusBadRequest)
		return
	}
	defer res.Body.Close()
	ok := a.db.ClearLogs(component)
	if ok {
		rd.Text(rw, http.StatusOK, "clear success")
	} else {
		rd.Text(rw, http.StatusOK, "there is no logs")
	}
}

func (a *AppHandler) CheckLogs(rw http.ResponseWriter, r *http.Request) {
	var novaerr, heaterr, cindererr, neutronerr, keystoneerr, swifterr, agenterr, managemneterr int
	errors := data.ComponentError{}
	novaerr = a.db.GetError("nova")
	heaterr = a.db.GetError("heat")
	cindererr = a.db.GetError("cinder")
	neutronerr = a.db.GetError("neutron")
	keystoneerr = a.db.GetError("keystone")
	swifterr = a.db.GetError("swift")
	agenterr = a.db.GetError("agent")
	managemneterr = a.db.GetError("management")

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
	if agenterr != 0 {
		errors.Agent = true
	}
	if managemneterr != 0 {
		errors.Management = true
	}

	rd.JSON(rw, http.StatusOK, errors)
}

func NewUnscopedTokenReq() data.TokenRequest {
	var newReq data.TokenRequest
	newReq.Auth.Identity.Methods = append(newReq.Auth.Identity.Methods, "password")
	domain := &data.Domain{}
	domain.Name = "Default"
	newReq.Auth.Identity.Password.User.Domain = domain
	return newReq
}

func NewScopedTokenReq() data.TokenRequest {
	var newReq data.TokenRequest
	newReq.Auth.Identity.Methods = append(newReq.Auth.Identity.Methods, "password")
	return newReq
}

func (a *AppHandler) GetToken(rw http.ResponseWriter, r *http.Request) {
	var login data.Login
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&login)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	UnscReq := NewUnscopedTokenReq()
	name := login.Name
	password := login.Password
	projectId := login.ProjectId
	UnscReq.Auth.Identity.Password.User.Name = name
	UnscReq.Auth.Identity.Password.User.Password = password

	rbytes, err := json.Marshal(UnscReq)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	buff := bytes.NewBuffer(rbytes)
	res, err := http.Post("http://164.125.70.22/identity/v3/auth/tokens", "application/json", buff)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	defer res.Body.Close()

	unscTokenBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	var unscToken map[string]map[string]map[string]string
	json.Unmarshal(unscTokenBody, &unscToken)
	id := unscToken["token"]["user"]["id"]

	ScReq := NewScopedTokenReq()
	ScReq.Auth.Identity.Password.User.Id = id
	ScReq.Auth.Identity.Password.User.Password = password
	scope := &data.Scope{}
	scope.Project.Id = projectId
	ScReq.Auth.Scope = scope

	rbytes, err = json.Marshal(ScReq)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	buff = bytes.NewBuffer(rbytes)

	res, err = http.Post("http://164.125.70.22/identity/v3/auth/tokens", "application/json", buff)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	defer res.Body.Close()

	token := res.Header.Get("X-Subject-Token")
	fmt.Println("token:", token)
	// resBody, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to read body", http.StatusBadRequest)
	// }
	var tokenRes struct {
		Token string `json:"token"`
	}
	// json.Unmarshal(resBody, &scopes)
	tokenRes.Token = token
	rd.JSON(rw, http.StatusOK, tokenRes)
}

func (a *AppHandler) GetInstances(rw http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	req, err := http.NewRequest("GET", "http://164.125.70.22/compute/v2.1/servers", nil)
	if err != nil {
		http.Error(rw, "Unable to get block", http.StatusBadRequest)
	}
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(rw, "Unable to do request", http.StatusInternalServerError)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(rw, "Unable to read body", http.StatusBadRequest)
	}
	var instances map[string]interface{}
	json.Unmarshal(resBody, &instances)
	rd.JSON(rw, http.StatusOK, instances)
}

func (a *AppHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	var multiMetrics data.Metrics

	token := r.Header.Get("token")

	req, err := http.NewRequest("GET", "http://164.125.70.22/compute/v2.1/os-hypervisors/statistics", nil)
	if err != nil {
		http.Error(rw, "Unable to get openstack statistics", http.StatusBadRequest)
	}
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(rw, "Unable to do request", http.StatusBadRequest)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(rw, "Unable to read body", http.StatusBadRequest)
	}
	var myHypervisor data.Hypervisors
	json.Unmarshal(resBody, &myHypervisor)
	fmt.Println(myHypervisor)

	res, err = http.Post("https://tyfgmh9pg3.execute-api.ap-northeast-2.amazonaws.com/PROD/signature", "application/json", r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	defer res.Body.Close()

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(rw, "Unable to read body", http.StatusBadRequest)
	}

	requestURL := string(resBody)
	requestURL = strings.Trim(requestURL, "\"")
	fmt.Println("http://164.125.70.26:8080/client/api?" + requestURL)

	req, err = http.NewRequest("GET", "http://164.125.70.26:8080/client/api?"+requestURL, nil)
	if err != nil {
		http.Error(rw, "Unable to get cloudstack statistics", http.StatusBadRequest)
	}
	req.Header.Set("content-type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		http.Error(rw, "Unable to do request", http.StatusBadRequest)
	}
	defer res.Body.Close()
	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(rw, "Unable to read body", http.StatusBadRequest)
	}

	var cloudstackMetrics data.CloudStackMetrics
	json.Unmarshal(resBody, &cloudstackMetrics)

	multiMetrics.OpenStackMetrics = myHypervisor.Statistics
	multiMetrics.CloudStackMetrics = cloudstackMetrics.Response

	rd.JSON(rw, http.StatusOK, multiMetrics)
}

func MakeHandler() *AppHandler {
	r := mux.NewRouter()
	neg := negroni.Classic()
	neg.UseHandler(r)

	a := &AppHandler{
		Handler: neg,
		db:      model.NewDBHandler(),
	}

	r.HandleFunc("/", a.IndexHandler)

	// Logging Handlers
	r.HandleFunc("/sync", a.SyncLogs).Methods("GET")
	r.HandleFunc("/{component:[a-z]+}/log", a.GetLogs).Methods("GET")
	r.HandleFunc("/{component:[a-z]+}/log", a.ClearLogs).Methods("DELETE")
	r.HandleFunc("/check", a.CheckLogs).Methods("GET")

	// Monitoring Handlers
	r.HandleFunc("/token", a.GetToken).Methods("POST")
	r.HandleFunc("/instances", a.GetInstances).Methods("GET")
	r.HandleFunc("/statistics", a.GetStatistics).Methods("POST")

	// Swagger Handlers
	opts := middleware.RedocOpts{SpecURL: "/swagger.yml"}
	sh := middleware.Redoc(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yml", http.FileServer(http.Dir("./")))

	return a
}
