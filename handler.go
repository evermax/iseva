package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

type JSONHandler struct {
	DB       string
	IsStatic bool
	dbc      dbContent
}

func (handler *JSONHandler) Init() error {
	return handler.getDBData()
}

func (handler *JSONHandler) getDBData() error {
	body, err := ioutil.ReadFile(handler.DB)
	if err != nil {
		return err
	}
	if strings.Contains(string(body), "---") {
		var params parameters
		bodyArr := strings.Split(string(body), "---")
		if err := json.Unmarshal([]byte(bodyArr[1]), &params); err != nil {
			return err
		}
		elts := params.parse()
		tmpl, err := template.New("JSONtemplate").Funcs(elts.Func).Parse(bodyArr[0])
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, elts.Var)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(buf.Bytes(), &handler.dbc); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(body, &handler.dbc); err != nil {
			return err
		}
	}
	return nil
}

func (handler *JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// support for cross domain options calls
	w.Header().Add("Access-Control-Allow-Header", "Content-Type, X-Requested-With")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if r.Method == "OPTIONS" {
		if origin := r.Header.Get("origin"); origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if !handler.IsStatic {
		err := handler.getDBData()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if raw, ok := handler.dbc.URLs[r.URL.Path]; ok {
		if origin := r.Header.Get("origin"); origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(raw.JSON)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

type dbContent struct {
	URLs map[string]raw `json:"urls"`
}

type raw struct {
	JSON json.RawMessage `json:"json"`
}

type parameters struct {
	Variables map[string]json.RawMessage `json:"variables"`
	Functions *funcParams                `json:"functions"`
}

func (params parameters) parse() tmplParams {
	var functions map[string]interface{}
	if params.Functions != nil {
		functions = params.Functions.parse()
	}
	var variables = make(map[string]string)
	for k, v := range params.Variables {
		variables[k] = string(v)
	}
	return tmplParams{
		Var:  variables,
		Func: functions,
	}
}

type tmplParams struct {
	Var  map[string]string
	Func map[string]interface{}
}
