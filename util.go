package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func getRecentFile() string {
	// get files in folder
	files, err := ioutil.ReadDir("E:\\replays")
	if err != nil {
		log.Fatal(err)
	}

	// get actual replays in directory
	replays := []string{""}
	for _, f := range files {
		if strings.Contains(f.Name(), "Replay") {
			replays = append(replays, f.Name())
		}
	}

	// get most recent replay by default
	lastReplay := replays[len(replays)-1]
	//log.Println("most recent replay:", lastReplay)

	return lastReplay
}

// is the file under 8MB? todo: accept param as flag, and/or set default in config file
func isFileTooBig(file string) bool {
	// stat file
	f, err := os.Stat(file)
	if err != nil {
		log.Println("problem getting file size", err)
		log.Fatal()
	}

	// get size
	size := f.Size()
	if size >= 8000000 {
		return true
	}

	return false
}
