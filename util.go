package main

import (
	"log"
	"os"
	"path/filepath"
)

// gets most recent file
// todo: change this from assuming ReadDir will list the most recent last to looking at creation timestamps, excluding directories
func getRecentFile() string {
	// get files in folder
	files, err := os.ReadDir(filepath.FromSlash(cfg.InputFolder))
	if err != nil {
		log.Fatal("problem getting most recent file:", err)
	}

	// get actual replays in directory
	replays := []string{""}
	for _, f := range files {
		if f.IsDir() == false {
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
	if size >= cfg.TargetFileSize {
		return true
	}

	return false
}
