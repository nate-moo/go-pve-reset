package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/reset", serveFile)
	http.HandleFunc("/reset/", resetHandler)
	http.HandleFunc("/reset/git-update", puller)

	http.HandleFunc("/manage", serveManage)
	http.HandleFunc("/manage/create", manage)
	http.HandleFunc("/manage/list.json", serveList)
	http.HandleFunc("/manage/manage.js", serveJS)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/reset.html")
}

func serveJS(writer http.ResponseWriter, request *http.Request) {
	b, err := os.ReadFile("static/manage.js")
	if err != nil {
		return
	}

	writer.Header().Add("Content-Type", "application/javascript")
	writer.Write(b)
}
