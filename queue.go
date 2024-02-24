package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Function to traverse directories and files
func addToQueue(rootPath string) {
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, add only file paths
		if !info.IsDir() {
			if !showIfav1(path) {
				queue = append(queue, path)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error traversing files: %v", err)
	}
}

func showIfav1(file string) bool {
	info, err := getMediaInfo(file)
	if err != nil {
		log.Fatalf("Error getting media info: %v", err)
	}
	return isAv1(info)
}

// func runqueue(concurrentNumber int) {

// }
// this is a little complex
func runQueue(limit int) {
	if len(queue) == 0 {
		fmt.Println("Queue is empty!")
		return
	}

	var wg sync.WaitGroup
	ch := make(chan struct{}, limit)

	for len(queue) > 0 {
		wg.Add(1)

		// Get the next file path from the queue
		filePath := queue[0]
		dirPath := filepath.Dir(filePath) // Get the directory of the original file
		cleanedName := cleanFileName(filepath.Base(filePath))
		newname := filepath.Join(dirPath, cleanedName+"-AV1.mp4") // Join the directory path with the new file name

		// Update the queue
		queue = queue[1:]

		ch <- struct{}{} // Block if limit is reached
		go func(file, newName string) {
			defer wg.Done()
			defer func() { <-ch }() // Release one spot
			err := encodeAv1(file, newName)
			if err != nil {
				fmt.Printf("Error encoding file: %v\n", err)
			} else {
				// Delete the original file after successful encoding
				err := os.Remove(file)
				if err != nil {
					fmt.Printf("Error deleting file: %v\n", err)
				} else {
					fmt.Printf("Deleted original file: %s\n", file)
				}
			}
		}(filePath, newname)
	}

	wg.Wait() // Wait for all goroutines to finish

	if len(queue) == 0 {
		fmt.Println("All tasks completed, queue is empty!")
	}
}

func cleanFileName(fileName string) string {
	// List of words to remove from the file name
	wordsToRemove := []string{"hevc", "h264", "h.264", "256", "HEVC", "x265"}
	// Compile a regex that matches any of the words in the list
	removeRegex := regexp.MustCompile(strings.Join(wordsToRemove, "|"))

	// Remove the words
	cleanedName := removeRegex.ReplaceAllString(fileName, "")

	// Remove the file extension
	cleanedName = strings.TrimSuffix(cleanedName, filepath.Ext(cleanedName))

	return cleanedName
}
