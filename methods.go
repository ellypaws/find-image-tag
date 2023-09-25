package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (data *DataSet) MoveCaptionsToImages() {
	// TODO: Move captions to images
	log.Print("Not implemented yet")
}

func (data *DataSet) CheckIfCaptionsExist() {
	for _, image := range data.Images {
		if image.Caption.Filename != "" {
			continue
		}
		fmt.Println("Caption for image", image.Filename, "does not exist")
	}
}

func (data *DataSet) WriteFiles() {
	regex := `(?i)\.(jpe?g|png|gif|bmp|txt)$`
	var directoryToRead string
	fmt.Print("Enter the directory to read: ")
	_, _ = fmt.Scanln(&directoryToRead)
	if directoryToRead == "" {
		directoryToRead = "."
	}

	extRegex, _ := regexp.Compile(regex)

	var tempCaption []Caption

	err := filepath.Walk(directoryToRead, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		currentEntry := info.Name()
		extension := extRegex.FindString(currentEntry)
		fileName := currentEntry[:len(currentEntry)-len(extension)]
		if !extRegex.MatchString(currentEntry) {
			return nil
		}

		directory, _ := filepath.Split(path)

		if extension == ".txt" {
			tempCaption = append(tempCaption, Caption{Filename: currentEntry, Extension: extension, Directory: directory})
		}

		data.Images[fileName] = Image{Filename: currentEntry, Extension: extension, Directory: directory, Caption: Caption{}}

		fmt.Println("Added file:", currentEntry, "to the dataset")

		return nil
	})
	if err != nil {
		return
	}

	data.appendCaptions(tempCaption)
}

func (data *DataSet) appendCaptions(c []Caption) {
	for _, caption := range c {
		fileName, _ := strings.CutSuffix(caption.Filename, caption.Extension)
		if img, ok := data.Images[fileName]; ok {
			fmt.Println("Appending the caption file:", caption.Filename, "to the image file:", fileName)
			img.Caption = caption
			data.Images[fileName] = img
		} else {
			fmt.Println("Image file for caption", caption.Filename, "does not exist")
		}
	}
}
