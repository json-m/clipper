package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

// config file format
type Config struct {
	InputFolder      string `yaml:"inputFolder"`
	OutputFolder     string `yaml:"outputFolder"`
	StartQuality     int    `yaml:"startQuality"`
	TargetResolution string `yaml:"targetResolution"`
	TargetCodec      string `yaml:"targetCodec"`
	TargetFileSize   int64  `yaml:"targetFileSize"`
	Audio            bool   `yaml:"audio"`
}

var cfg Config

func init() {
	// test for ffmpeg
	fInstalled, fCompatible := testFfmpeg()

	// bail if ffmpeg isn't found/compatible
	if !fInstalled || !fCompatible {
		if !fInstalled {
			log.Fatal("ffmpeg not found")
		}
		if !fCompatible {
			log.Fatal("ffmpeg doesn't support nvenc")
		}
	}

	// open config file
	f, err := os.Open("clipper.yml")
	if err != nil {
		log.Fatal("couldn't open config file:", err)
	}
	defer f.Close()

	// read config file
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("couldn't read config:", err)
	}

	// set log output
	log.SetOutput(os.Stdout)

	// are the minimum count of args set?
	if len(os.Args) >= 3 {
		//...
	} else {
		log.Fatal(errors.New("syntax error: too few args, syntax is clipper.exe -time=##:## -dur=# [-aud=bool] [-file=replay.mkv]"))
	}

}

func main() {
	// read flags
	var timeArg string
	var fileArg string
	flag.StringVar(&timeArg, "time", "00:00", "-time=02:35")
	durationArg := flag.Int("dur", 00, "-dur=12")
	audioArg := flag.Bool("aud", cfg.Audio, "-aud=true")
	flag.StringVar(&fileArg, "file", "", "-file=replay.mkv")
	flag.Parse()

	// another simple check for args here
	if timeArg == "00:00" && *durationArg == 00 {
		log.Fatal(errors.New("syntax error: -time and -dur left empty"))
	}

	// set inputFile
	var inputFile string
	if fileArg != "" {
		// does the desired input file even exist?
		if _, err := os.Stat(fileArg); os.IsNotExist(err) {
			log.Fatal("input file does not exist:", err)
		}
		inputFile = fileArg
		log.Println("input file was specified:", inputFile)
	} else {
		inputFile = getRecentFile()
		log.Println("most recent replay:", inputFile)
	}

	// filename based on unixtime
	outputFilename := fmt.Sprintf("%s\\%d.mp4", filepath.FromSlash(cfg.OutputFolder), time.Now().Unix())

	// clip the video using ffmpeg
	err := ffmpegClip(timeArg, *durationArg, cfg.StartQuality, inputFile, outputFilename, *audioArg)
	if err != nil {
		log.Fatal("problem calling ffmpegClip", err)
	}

	// check the output file size - if it's greater than desired, increase QP
	for {
		if isFileTooBig(outputFilename) {
			cfg.StartQuality++
			log.Println("file too big, increased QP:", cfg.StartQuality)
			err := ffmpegClip(timeArg, *durationArg, cfg.StartQuality, inputFile, outputFilename, *audioArg)
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
