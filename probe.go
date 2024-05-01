package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// Define the struct to match the JSON structure
type MediaData struct {
	Streams []struct {
		Index int `json:"index"`
		Tags  struct {
			Language    string `json:"language"`
			HandlerName string `json:"handler_name"`
			VendorID    string `json:"vendor_id"`
		} `json:"tags"`
	} `json:"streams"`
	Format struct {
		Filename string `json:"filename"`
	} `json:"format"`
}

func isAv1(probe string) bool {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(probe), &data); err != nil {
		log.Fatal(err)
		return false
	}
	streams, ok := data["streams"].([]interface{})
	if !ok || len(streams) == 0 {
		fmt.Println("No streams found or unable to parse streams.")
		return false
	}

	firstStream, ok := streams[0].(map[string]interface{})
	if !ok {
		fmt.Println("Unable to parse the first stream.")
		return false
	}

	codecName, ok := firstStream["codec_name"].(string)
	if !ok {
		fmt.Println("Codec name not found or is not a string.")
		return false
	}
	filename, _ := getFileName(probe)
	if filename == "" {
		fmt.Println("cant get filename for below")
	}
	fmt.Println("Codec Name:", codecName+" <- "+filename)
	return codecName == av1
}

func getMediaInfo(file string) (string, error) {
	// Constructing the ffprobe command
	// This command will output information in a user-friendly format
	cmd := exec.Command("ffprobe", "-v", "error", "-show_format", "-show_streams", "-print_format", "json", file)

	// Capturing the output
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("() getMediaInfo : ffprobe error: %w", err)
	}
	return out.String(), nil
}

func getFileName(path string) (string, error) {
	var mediaData MediaData
	err := json.Unmarshal([]byte(path), &mediaData)
	if err != nil {
		return "", err
	}
	return mediaData.Format.Filename, nil
}
