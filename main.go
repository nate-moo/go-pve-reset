package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func serveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/root/webserver/pve/test.html")
}

func main() {
	http.HandleFunc("/reset", serveFile)
	http.HandleFunc("/reset/", resetHandler)
	http.HandleFunc("/reset/git-update", puller)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

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

func resetHandler(writer http.ResponseWriter, request *http.Request) {
	split := strings.Split(request.URL.String(), "/")

	type vmFormat struct {
		templateID int
		vmID       [10]int
		name       string
	}

	// Templates
	var (
		UbuntuEasy1      = 110
		UbuntuEasy2      = 111
		UbuntuMedium1    = 120
		UbuntuMedium2    = 121
		SupportEasy1     = 210
		SupportEasy2     = 211
		UbuntuPlayground = 500
	)

	// mappings
	var VMMap = map[int]vmFormat{
		UbuntuEasy1:      {templateID: 110, vmID: [10]int{1101}, name: "Ubuntu-Easy-1"},
		UbuntuEasy2:      {templateID: 111, vmID: [10]int{1111}, name: "Ubuntu-Easy-2"},
		UbuntuMedium1:    {templateID: 120, vmID: [10]int{1201}, name: "Ubuntu-Medium-1"},
		UbuntuMedium2:    {templateID: 121, vmID: [10]int{1211}, name: "Ubuntu-Medium-2"},
		SupportEasy1:     {templateID: 210, vmID: [10]int{2101}, name: "Support-Easy-1"},
		SupportEasy2:     {templateID: 211, vmID: [10]int{2111}, name: "Support-Easy-2"},
		UbuntuPlayground: {templateID: 500, vmID: [10]int{5000}, name: "Ubuntu-Playground"},
	}

	var skip = false

	targetName := split[4]

	templateID, err := strconv.Atoi(split[2])
	if err != nil {
		skip = true
		log.Println("Failed to convert")
	}

	targetID, err := strconv.Atoi(split[3])
	if err != nil {
		skip = true
		log.Println("Failed to convert")
	}

	var corr = false

	if (targetID / 10) != templateID {
		skip = true
	}

	if !skip {
		for _, vm := range VMMap {
			if vm.templateID == templateID {
				for _, vmID := range vm.vmID {
					if vmID == targetID && vm.name == targetName {
						corr = true
					}
				}
			}
		}
	}

	if corr == false || skip {
		errorMSG := "Request Parameter Match Error. Check for valid VMID, TemplateID, or Name"
		writer.WriteHeader(400)

		_, err := writer.Write([]byte(errorMSG))
		if err != nil {
			return
		}

		log.Fatal(errorMSG)
		return
	}

	err = exec.Command("qm", "stop", strconv.Itoa(targetID)).Run()
	if err != nil {
		log.Println("qm stop " + strconv.Itoa(targetID) + "failed")
	}

	err = exec.Command("qm", "destroy", strconv.Itoa(targetID)).Run()
	if err != nil {
		log.Println("qm destroy " + strconv.Itoa(targetID) + "failed")
	}

	err = exec.Command("qm", "clone", strconv.Itoa(templateID), strconv.Itoa(targetID), "--name", targetName).Run()
	if err != nil {
		log.Println("qm clone " + strconv.Itoa(templateID) + "->" + strconv.Itoa(targetID) + "failed")
	}

	// pvesh set /access/acl -path /vms/{vmid} -roles PVEVMUser -groups SD_Users
	cmd := exec.Command("pvesh", "set", "/access/acl", "-path", "/vms/"+strconv.Itoa(targetID), "-roles", "PVEVMUser", "-groups", "SD_Users")

	err = cmd.Run()
	if err != nil {
		log.Println("pvesh " + "set " + "/access/acl " + "-path " + "/vms/" + strconv.Itoa(targetID) + " -roles " + "PVEVMUser " + "-groups " + "SD_Users " + "failed")
	}

	log.Println("Template: "+strconv.Itoa(templateID), "Target: "+strconv.Itoa(targetID), targetName)

	writer.WriteHeader(200)
	msg := fmt.Sprintf("<html><head><title>Reset Success!</title></head><body>Reset of VM %s (%d -> %d) Successful!</body></html>", targetName, templateID, targetID)
	log.Println(msg)
	_, err = writer.Write([]byte(msg))
	if err != nil {
		return
	}
}
