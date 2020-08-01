package main

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Game struct {
	File  string `json:"file"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Games []Game

func New(file string) Game {
	s := strings.Split(file, ".")
	ng := Game{}
	ng.File = file
	ng.Name = s[0]
	ng.Image = gameImage(ng)
	return ng
}

func gameImage(game Game) string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	rd := path + "/assets/" + game.Name + ".jpg"
	info, err := os.Stat(rd)
	if os.IsNotExist(err) {
		return defaultImage()
	}
	if info.IsDir() {
		return defaultImage()
	}
	return readImageFile(rd)
}

func defaultImage() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	rd := path + "/assets/no_covers.jpg"
	return readImageFile(rd)
}

func readImageFile(file string) string {
	f, _ := os.Open(file)

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded
}
