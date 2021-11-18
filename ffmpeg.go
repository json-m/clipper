package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func ffmpegClip(startTime string, durationSec, qp int, inputFile, outputFile string, audio bool) error {
	var err error

	// start building ffmpeg args
	args := []string{"-y", "-hide_banner", // yes to overwriting, hide the giant banner of supported things
		"-ss", fmt.Sprintf("00:%s", startTime), // start timestamp
		"-i", fmt.Sprintf("%s\\%s", cfg.InputFolder, inputFile), // input file
		"-t", fmt.Sprintf("00:00:%d", durationSec), // length of clip
		"-vcodec", "h264_nvenc", // nvidia encoder
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	runFfmpeg.Stdout = &stdout
	runFfmpeg.Stderr = &stderr

	// now run ffmpeg
	err = runFfmpeg.Run()
	if err != nil {
		log.Println("ffmpeg command was:", runFfmpeg.Args)
		log.Println("failed to run", stderr.String())
		log.Fatal(err)
	}

	return err
}
