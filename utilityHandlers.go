package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

func puller(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}

	sec := request.Form.Get("security")

	file, err := os.Open("/root/webserver/secret")
	if err != nil {
		return
	}

	secCheck := []byte("1234567890")
	_, err = file.Read(secCheck)
	if err != nil {
		return
	}

	if sec == string(secCheck) {
		err := os.Chdir("/root/webserver/pve")
		if err != nil {
			log.Println("can't chdir")
			return
		}

		log.Println("Success")

		err = exec.Command("git", "pull").Run()
		if err != nil {
			log.Println("git pull failed")
			return
		}

		err = exec.Command("go", "build", ".").Run()
		if err != nil {
			log.Println("go build failed")
			return
		}

		err = exec.Command("systemctl", "restart", "pve-reset.service").Run()
		if err != nil {
			log.Println("service restart failed")
			return
		}
	}
}
