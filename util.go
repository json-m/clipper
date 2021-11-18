package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// gets most recent file
// todo: change this from assuming ReadDir will list the most recent last to looking at creation timestamps, excluding directories
// todo: "As of Go 1.16, os.ReadDir is a more efficient and correct choice" - https://pkg.go.dev/io/ioutil#ReadDir
func getRecentFile() string {
	// get files in folder
	files, err := ioutil.ReadDir(cfg.InputFolder)
	if err != nil {
		log.Fatal("problem getting most recent file:", err)
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

	return lastReplay
}

// is the file under the configured size?
func isFileTooBig(file string) bool {
	// stat file
	f, err := os.Stat(file)
	if err != nil {
		log.Fatal("problem getting file size:", err)
	}

	// get size
	size := f.Size()
	if size >= cfg.TargetFileSize { // 8MB = non-nitro discord
		return true
	}

	return false
}
