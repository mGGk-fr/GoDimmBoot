package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var g Game.Games
var rd string

func main() {
	fmt.Println("---------------")
	fmt.Println("- Netboot Web -")
	fmt.Println("- Version 0.1 -")
	fmt.Println("---------------")
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	rd = path + "/roms"
	fmt.Println("Roms directory : " + rd)
	fmt.Println("Please wait, scanning games...")
	scanRoms()
	fmt.Println("Game scan ended")
	fmt.Println("Booting web server...")
	bootWebServer()
}

func scanRoms() {
	g = nil
	libRegEx, e := regexp.Compile("^.+\\.(bin)$")
	if e != nil {
		log.Fatal(e)
	}

	e = filepath.Walk(rd, func(p string, i os.FileInfo, e error) error {
		if e == nil && libRegEx.MatchString(i.Name()) {
			fmt.Println("Found " + i.Name())
			g = append(g, New(i.Name()))
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}
}

func bootWebServer() {
	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", fs)
	http.HandleFunc("/reload", reloadGameList)
	http.HandleFunc("/games", sendGameList)
	http.HandleFunc("/sendGame", sendGame)
	http.ListenAndServe(":8080", nil)
}

func sendGameList(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(g)
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Fprint(w, string(b))
}

func reloadGameList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rescan ")
	scanRoms()
	b, _ := json.Marshal(g)
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Fprint(w, string(b))
}

func sendGame(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["game"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	key := keys[0]

	log.Println("Url Param 'key' is: " + string(key))
	cmd := exec.Command("python", "support/booter.py", key)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(string(out))
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
