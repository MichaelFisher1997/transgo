package main

import (
	"os/exec"
)

// convert file to av1 using av1_qsv codec and hwaccel vaapi and hwaccel_output_format vaapi
func encodeAv1(file string, newname string) error {
	// Construct the FFmpeg command with hardware acceleration and codec settings
	cmd := exec.Command("ffmpeg", "-hwaccel", "vaapi", "-hwaccel_output_format", "vaapi",
		"-i", file, "-vf", "hwmap=derive_device=qsv,format=qsv", "-c:v", "av1_qsv", "-preset", "slow", newname)

	// Execute the command
	err := cmd.Run()
	if err != nil {
		//fmt.Printf("FFmpeg error: %v\n", err)
		return err
	}

	//log.Println("Video encoding completed successfully")
	return nil
}
