package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (data *DataSet) MoveCaptionsToImages() {
	// TODO: Move captions to images
	roggyPrinter.Debug("Not implemented yet")
}

func (data *DataSet) CheckIfCaptionsExist() {
	for _, image := range data.Images {
		if image.Caption.Filename != "" {
			continue
		}
		roggyPrinter.Errorf("Caption for image %s does not exist", image.Filename)
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

		roggyPrinter.Debugf("Filename: %s", fileName)
		roggyPrinter.Debugf("Extension: %s", extension)

		if extension == ".txt" {
			tempCaption = append(tempCaption, Caption{Filename: currentEntry, Extension: extension, Directory: directory})
			return nil
		}

		if _, ok := data.Images[fileName]; !ok {
			data.Images[fileName] = Image{Filename: currentEntry, Extension: extension, Directory: directory, Caption: Caption{}}
			roggyPrinter.Infof("Added file: %s to the dataset", currentEntry)
			roggyPrinter.Debugf("Directory: %s", directory)
		}

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
			roggyPrinter.Infof("Appending the caption file: %s to the image file: %s", caption.Filename, fileName)
			roggyPrinter.Debugf("Directory: %s", caption.Directory)
			img.Caption = caption
			data.Images[fileName] = img
		} else {
			roggyPrinter.Noticef("Image file for caption %s does not exist", caption.Filename)
		}
	}
}
