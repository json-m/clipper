package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// check for ffmpeg
func testFfmpeg() (bool, bool) {
	var isInstalled bool
	var hasNVIDIA bool

	// tests that ffmpeg is even installed
	pathTest := exec.Command("ffmpeg", "-L")
	err := pathTest.Run()
	if err != nil {
		return false, false
	}

	// didn't get error while executing
	isInstalled = true

	// setup test for nvenc encoders
	nvencTest := exec.Command("ffmpeg", "-hide_banner", "-encoders")

	// run ffmpeg cmd and reformat output text
	output, err := nvencTest.CombinedOutput()
	if err != nil {
		return false, false
	}
	out := strings.Split(string(output), "\n")

	// check output for nvenc encoder
	for _, line := range out {
		if strings.Contains(line, "nvenc") {
			hasNVIDIA = true
			break
		}
	}

	return isInstalled, hasNVIDIA
}

func ffmpegClip(startTime string, durationSec, qp int, inputFile, outputFile string, audio bool) error {
	// start building ffmpeg args
	args := []string{"-y", "-hide_banner", // yes to overwriting, hide the giant banner of supported things
		"-ss", fmt.Sprintf("00:%s", startTime), // start timestamp
		"-i", fmt.Sprintf("%s\\%s", cfg.InputFolder, inputFile), // input file
		"-t", fmt.Sprintf("00:00:%d", durationSec), // length of clip
		"-vcodec", cfg.TargetCodec, // nvidia encoder
		"-s", cfg.TargetResolution, // output resolution
		"-rc", "constqp", // i know why constqp is *typically* bad, and in this specific case, i do not care
		"-qp", fmt.Sprintf("%d", qp), // tune for nvenc constqp
	} // https://archive.md/8YzUL - "Understanding Rate Control Modes (x264, x265, vpx)"

	// if audio is desired set some opts for it. otherwise, drop the audio track
	if audio {
		args = append(args, "-b:a")
		args = append(args, "128k")
	} else {
		args = append(args, "-an")
	}

	// finally, append the output file
	args = append(args, fmt.Sprintf("%s", outputFile))

	// create the command
	runFfmpeg := exec.Command("ffmpeg", args...)

	// setup for output
	var stderr bytes.Buffer
	runFfmpeg.Stderr = &stderr

	// now run ffmpeg
	err := runFfmpeg.Run()
	if err != nil {
		log.Println("ffmpeg command was:", runFfmpeg.Args)
		log.Println("failed to run", stderr.String())
		return err
	}

	return nil
}
