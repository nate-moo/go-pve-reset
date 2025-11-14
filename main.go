package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func serveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "test.html")
}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/reset", serveFile)
	//http.HandleFunc("/test.html", serveFile)
	http.HandleFunc("/reset/", resetHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func resetHandler(writer http.ResponseWriter, request *http.Request) {
	split := strings.Split(request.URL.String(), "/")

	err := exec.Command("qm", "stop", split[3]).Run()
	if err != nil {
		log.Println("qm stop " + split[3] + "failed")
	}

	err = exec.Command("qm", "destroy", split[3]).Run()
	if err != nil {
		log.Println("qm destroy " + split[3] + "failed")
	}

	err = exec.Command("qm", "clone", split[2], split[3], "--name", split[4]).Run()
	if err != nil {
		log.Println("qm clone " + split[2] + "->" + split[3] + "failed")
	}

	// pvesh set /access/acl -path /vms/{vmid} -roles PVEVMUser -groups SD_Users
	cmd := exec.Command("pvesh", "set", "/access/acl", "-path", "/vms/"+split[3], "-roles", "PVEVMUser", "-groups", "SD_Users")

	err = cmd.Run()
	if err != nil {
		log.Println("pvesh " + "set " + "/access/acl " + "-path " + "/vms/" + split[3] + " -roles " + "PVEVMUser " + "-groups " + "SD_Users " + "failed")
	}

	log.Println(split[2], split[3], split[4])
}
