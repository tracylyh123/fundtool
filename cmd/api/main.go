package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tracylyh123/fundtool/fund"
)

type ret struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
}

func (r ret) toJSON() []byte {
	s, _ := json.Marshal(r)
	return s
}

const (
	failed  = "failed"
	success = "success"
)

func main() {
	http.HandleFunc("/trend", trend)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func trend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		code = r.URL.Query().Get("code")
		from = r.URL.Query().Get("from")
		to   = r.URL.Query().Get("to")
	)
	var status int = http.StatusOK
	var ret ret
	if len(code) == 0 {
		ret.Status = failed
		ret.Msg = "code was empty"
		status = http.StatusBadRequest
	} else {
		trend, err := fund.GetFundTrend(code, from, to)
		if err == nil {
			ret.Status = success
			ret.Data = trend
		} else {
			ret.Status = failed
			ret.Msg = err.Error()
			status = http.StatusInternalServerError
		}
	}
	w.WriteHeader(status)
	w.Write(ret.toJSON())
}
