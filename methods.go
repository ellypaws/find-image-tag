package main

import (
	"fmt"
	"github.com/nokusukun/roggy"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var captionLogPrinter = roggy.Printer("caption-handler")

const (
	move = iota
	hardlink
	merge
)

func (data *DataSet) checkForMissingImages() {
	wg := &sync.WaitGroup{}

	for index := range data.TempCaption {
		wg.Add(1)
		go func(c *Caption) {
			defer wg.Done()
			tempCaptionLogPrinter := roggy.Printer(fmt.Sprintf("caption-handler: %v", c.Filename))
			fileName, _ := strings.CutSuffix(c.Filename, c.Extension)
			if _, ok := data.Images[fileName]; !ok {
				tempCaptionLogPrinter.Errorf("Image file for caption %s does not exist", c.Filename)
			}
		}(data.TempCaption[index])
	}

	wg.Wait()
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

		fileLogPrinter := roggy.Printer(fmt.Sprintf("file-handler: %v", fileName))
		tempCaptionLogPrinter := roggy.Printer(fmt.Sprintf("temp-caption-handler: %v", currentEntry))
		tempCaptionLogPrinter.NoTrace = true
		imageLogPrinter := roggy.Printer(fmt.Sprintf("image-handler: %v", currentEntry))

		directory, _ := filepath.Split(path)

		fileLogPrinter.Debugf("Processing File: %s", currentEntry)

		if extension == ".txt" {
			if _, ok := data.TempCaption[fileName]; !ok {
				if data.TempCaption == nil {
					data.InitTempCaption()
				}
				newCaption := Caption{Filename: currentEntry, Extension: extension, Directory: directory}
				data.TempCaption[fileName] = &newCaption
				tempCaptionLogPrinter.Infof("Added file: %s to the temporary caption dataset", currentEntry)
				tempCaptionLogPrinter.Debugf("Directory: %s", directory)
			} else {
				tempCaptionLogPrinter.Errorf("Caption file in temporary caption dataset for image %s already exists", fileName)
			}
			return nil
		}

		if _, ok := data.Images[fileName]; !ok {
			newImg := Image{Filename: currentEntry, Extension: extension, Directory: directory}
			newImg.InitCaption()
			data.Images[fileName] = &newImg
			imageLogPrinter.Infof("Added file: %s to the dataset", currentEntry)
			imageLogPrinter.Debugf("Directory: %s", directory)
		}

		return nil
	})
	if err != nil {
		return
	}

	captionLogPrinter.Infof("Now appending captions to images...")
	data.appendCaptionsConcurrently()
}
