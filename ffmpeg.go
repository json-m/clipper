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
	args := []string{"-y",
		"-ss", fmt.Sprintf("00:%s", startTime), // start timestamp
		"-i", fmt.Sprintf("E:\\replays\\%s", inputFile), // input file
		"-t", fmt.Sprintf("00:00:%d", durationSec), // length of clip
		"-vcodec", "h264_nvenc", // nvidia encoder
		"-s", "1280x720", // resolution to 720p
		"-rc", "constqp", // tune for nvenc
		"-qp", fmt.Sprintf("%d", qp), // tune for nvenc (basically CRF)
	}

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
