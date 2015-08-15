package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Interface struct {
	DB string
}

func (env *Interface) getDBData() (*dbContent, error) {
	body, err := ioutil.ReadFile(env.DB)
	if err != nil {
		return nil, err
	}
	var jsf dbContent
	if err := json.Unmarshal(body, &jsf); err != nil {
		return nil, err
	}
	return &jsf, nil
}

func (env *Interface) Handler(w http.ResponseWriter, r *http.Request) {
	// support for cross domain options calls
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if r.Method == "OPTIONS" {
		if origin := r.Header.Get("origin"); origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		w.Header().Add("Access-Control-Allow-Header", "Content-Type, X-Requested-With")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	db, err := env.getDBData()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if raw, ok := db.URLs[r.URL.Path]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		if origin := r.Header.Get("origin"); origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(raw.JSON)
	}
}

type dbContent struct {
	URLs      map[string]raw    `json:"URLs"`
	Variables map[string]string `json:"variables"`
}

type raw struct {
	JSON json.RawMessage `json:"JSON"`
}
