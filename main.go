package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func init() {
	// todo: does ffmpeg exist?

	// todo: get->set paths, perhaps from config file?

	// set log output
	log.SetOutput(os.Stdout)

	// are the 3 minimum flags set?
	if len(os.Args) >= 3 {
		//...
	} else {
		log.Fatal(errors.New("syntax error: too few args, syntax is clipper.exe -time=##:## -dur=# [-aud=bool] [-file=replay.mkv]"))
	}

}

func main() {
	// get most recent file
	inputFile := getRecentFile()

	// read flags
	var timeArg string
	var fileArg string
	flag.StringVar(&timeArg, "time", "00:00", "-time=02:35")
	durationArg := flag.Int("dur", 00, "-dur=12")
	audioArg := flag.Bool("aud", false, "-aud=true")
	flag.StringVar(&fileArg, "file", "", "-file=replay.mkv")
	flag.Parse()

	// another simple check for args here
	if timeArg == "00:00" && *durationArg == 00 {
		log.Fatal(errors.New("syntax error: -time and -dur left empty"))
	}

	// if input file was specified, update inputFile
	if fileArg != "" {
		inputFile = fileArg
		log.Println("input file was specified:", inputFile)
	} else {
		log.Println("most recent replay:", inputFile)
	}

	// todo: move to config file?
	nvencQp := 29 // seems like a good start

	// name based on unixtime
	outputFilename := fmt.Sprintf("%d.mp4", time.Now().Unix())

	// setup ffmpeg command
	err := ffmpegClip(timeArg, *durationArg, nvencQp, inputFile, outputFilename, *audioArg)
	if err != nil {
		log.Println("problem calling ffmpegClip", err)
		log.Fatal(err)
	}

	// check the output file size, if it's greater than 8MB (non-nitro discord), increase QP
	for {
		if isFileTooBig(outputFilename) {
			nvencQp++
			log.Println("file too big, increased QP:", nvencQp)
			err := ffmpegClip(timeArg, *durationArg, nvencQp, inputFile, outputFilename, *audioArg)
			if err != nil {
				log.Println("problem calling ffmpegClip", err)
				log.Fatal(err)
			}
		} else {
			// file is smaller than specified limit
			break
		}
	}

	log.Println("created clip:", outputFilename)
}
