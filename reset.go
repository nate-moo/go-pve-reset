package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func resetHandler(writer http.ResponseWriter, request *http.Request) {
	split := strings.Split(request.URL.String(), "/")

	b, err := os.ReadFile("vms.json")
	if err != nil {
		log.Println("Failed to read vms.json")
		return
	}

	var vmList []vmFormat

	err = json.Unmarshal(b, &vmList)
	if err != nil {
		log.Println("Failed to parse vms.json")
		return
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
		for _, vm := range vmList {
			if vm.TemplateID == templateID {
				for _, vmID := range vm.VmID {
					if vmID == targetID && vm.Name == targetName {
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

	update(templateID, targetID, targetName, true)

	log.Println("Template: "+strconv.Itoa(templateID), "Target: "+strconv.Itoa(targetID), targetName)

	writer.WriteHeader(200)
	msg := fmt.Sprintf("<html><head><title>Reset Success!</title></head><body>Reset of VM %s (%d -> %d) Successful!</body></html>", targetName, templateID, targetID)
	log.Println(msg)
	_, err = writer.Write([]byte(msg))
	if err != nil {
		return
	}
}

func update(templateID int, targetID int, targetName string, reset bool) {
	if reset {
		err := exec.Command("qm", "stop", strconv.Itoa(targetID)).Run()
		if err != nil {
			log.Println("qm stop " + strconv.Itoa(targetID) + "failed")
		}

		err = exec.Command("qm", "destroy", strconv.Itoa(targetID)).Run()
		if err != nil {
			log.Println("qm destroy " + strconv.Itoa(targetID) + "failed")
		}
	}

	err := exec.Command("qm", "clone", strconv.Itoa(templateID), strconv.Itoa(targetID), "--name", targetName).Run()
	if err != nil {
		log.Println("qm clone " + strconv.Itoa(templateID) + "->" + strconv.Itoa(targetID) + "failed")
	}

	// pvesh set /access/acl -path /vms/{vmid} -roles PVEVMUser -groups SD_Users
	cmd := exec.Command("pvesh", "set", "/access/acl", "-path", "/vms/"+strconv.Itoa(targetID), "-roles", "PVEVMUser", "-groups", "SD_Users")

	err = cmd.Run()
	if err != nil {
		log.Println("pvesh " + "set " + "/access/acl " + "-path " + "/vms/" + strconv.Itoa(targetID) + " -roles " + "PVEVMUser " + "-groups " + "SD_Users " + "failed")
	}
}
