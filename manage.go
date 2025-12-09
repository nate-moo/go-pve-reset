package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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

func manage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}
	vm := request.Form.Get("vm")
	TemplateID, err := strconv.Atoi(vm)

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

	nextID := -100
	cont := false
	var name string
	var index int
	var jindex int

	for i, v := range vmList {
		if v.TemplateID == TemplateID {
			name = v.Name
			index = i
			for j, id := range v.VmID {
				if id != 0 && id > nextID {
					cont = true
					jindex = j
					log.Println("found next highest ID, " + strconv.Itoa(int(nextID)))
				} else {
					log.Println("Something broke, got non-real vmid" + strconv.Itoa(nextID))
				}
			}
		}
	}

	// 10000 plasma
	// 00000 gnome -> 0000
	if cont {
		// == 0 means something broke, == 9 means we exhausted our vmid limit
		if nextID%10 != 0 || nextID%10 != 9 {
			nextID += 1
			update(TemplateID, nextID, name, false)
			vmList[index].VmID[jindex] = nextID
		} else {
			log.Println("Something broke: " + strconv.Itoa(nextID))
		}

	}

}

func serveManage(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("static/manage.html")
	if err != nil {
		return
	}

	w.Write(bytes)
}
