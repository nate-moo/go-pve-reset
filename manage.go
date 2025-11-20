package main

import (
	"net/http"
	"os"
)

func serveList(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("vms.json")
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func manage(w http.ResponseWriter, r *http.Request) {

}

func serveManage(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("static/manage.html")
	if err != nil {
		return
	}

	w.Write(bytes)
}
