package main

import (
	"bufio"
	"fmt"
	"github.com/nokusukun/roggy"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	move = iota
	hardlink
	merge
)

func (data *DataSet) CaptionsToImages(action int) {
	for _, image := range data.Images {
		if image.Caption.Filename == "" {
			captionLogPrinter.Noticef("Image %s does not have a caption", image.Filename, "Skipping...")
			continue
		}

		// Handle relative directories, get current directory with Get working directory
		currentDir, _ := os.Getwd()

		// If the directory is not an absolute path, join it with the current directory
		if !filepath.IsAbs(image.Caption.Directory) {
			image.Caption.Directory = filepath.Join(currentDir, image.Caption.Directory)
		}
		if !filepath.IsAbs(image.Directory) {
			image.Directory = filepath.Join(currentDir, image.Directory)
		}

		// Create the paths
		from := filepath.Join(image.Caption.Directory, image.Caption.Filename)
		to := filepath.Join(image.Directory, image.Caption.Filename)

		switch action {
		case move:
			// Move the file. On many file systems, this is a simple rename operation.
			err := os.Rename(from, to)
			if err != nil {
				captionLogPrinter.Errorf("Error moving file: %v", err)
			} else {
				captionLogPrinter.Noticef("File moved successfully from %s to %s", from, to)
			}
		case hardlink:
			// Create a hardlink of the file.
			err := os.Link(from, to)
			if err != nil {
				captionLogPrinter.Errorf("Error linking file: %v", err)
			} else {
				captionLogPrinter.Infof("File linked successfully from %s to %s", from, to)
			}
		case merge:
			// Read the contents of the existing caption
			existingContentBytes, err := os.ReadFile(from)
			if err != nil {
				captionLogPrinter.Errorf("Cannot read file: %s. Error: %v", from, err)
				continue
			}

			// Read the contents of the new caption
			newContentBytes, err := os.ReadFile(to)
			if err != nil {
				captionLogPrinter.Errorf("Cannot read file: %s. Error: %v", to, err)
				continue
			}

			existingContent := string(existingContentBytes)
			newContent := string(newContentBytes)

			// Split the contents by comma
			existingTags := strings.Split(existingContent, ",")
			newTags := strings.Split(newContent, ",")

			// Trim spaces
			for i, tag := range existingTags {
				existingTags[i] = strings.TrimSpace(tag)
			}
			for i, tag := range newTags {
				newTags[i] = strings.TrimSpace(tag)
			}

			// Combine the tags, eliminating duplicates
			tagMap := make(map[string]bool)
			for _, tag := range existingTags {
				tagMap[tag] = true
			}
			for _, tag := range newTags {
				tagMap[tag] = true
			}

			// Make a slice of the combined tags
			combinedTags := make([]string, 0, len(tagMap))
			for tag := range tagMap {
				combinedTags = append(combinedTags, tag)
			}

			// Join the tags with commas
			combinedContent := strings.Join(combinedTags, ", ")

			// Write the combined content to the new caption file
			err = os.WriteFile(to, []byte(combinedContent), 0644)
			if err != nil {
				captionLogPrinter.Errorf("Cannot write to file: %s. Error: %v", to, err)
			} else {
				captionLogPrinter.Infof("File combined successfully from %s to %s", from, to)
			}
			// Update the image.Caption.Directory to the new directory 'to'
			image.Caption.Directory = filepath.Dir(to)
		}
	}
}

func (data *DataSet) replaceSpaces() {
	re := regexp.MustCompile(`(\w)\s+(\w)`)
	wg := &sync.WaitGroup{}

	for index := range data.Images {
		wg.Add(1)
		go func(i *Image) {
			defer wg.Done()
			if i.Caption.Filename == "" {
				captionLogPrinter.Noticef("Image %s does not have a caption", i.Filename, "Skipping...")
				return
			}

			captionFile := filepath.Join(i.Caption.Directory, i.Caption.Filename)

			file, err := os.Open(captionFile)
			if err != nil {
				captionLogPrinter.Errorf("Error opening file: %v", err)
				return
			}

			reader := bufio.NewReader(file)
			content, _ := reader.ReadString('\n')
			_ = file.Close()

			newContent := re.ReplaceAllString(content, "${1}_${2}")

			err = os.Remove(captionFile)
			if err != nil {
				captionLogPrinter.Errorf("Error deleting file: %v", err)
				return
			}

			file, err = os.Create(captionFile)
			if err != nil {
				captionLogPrinter.Errorf("Error creating file: %v", err)
				return
			}

			writer := bufio.NewWriter(file)
			_, err = writer.WriteString(newContent + "\n")
			_ = writer.Flush()
			_ = file.Close()

			if err != nil {
				captionLogPrinter.Errorf("Error writing file: %v", err)
			} else {
				captionLogPrinter.Infof("Replaced spaces with underscores for file: %s", captionFile)
			}
		}(data.Images[index])
	}
	wg.Wait()
}

func (data *DataSet) CheckIfCaptionsExist() {
	var wg sync.WaitGroup
	for index := range data.Images {
		wg.Add(1)
		go func(image *Image) {
			if image.Caption.Filename != "" {
				roggyPrinter.Noticef("Caption `%s` has a matching image file `%s`", image.Caption.Filename, image.Filename)
				roggyPrinter.Debugf("Caption directory: %s", image.Caption.Directory)
				roggyPrinter.Debugf("Image directory: %s", image.Directory)
				return
			}
			roggyPrinter.Errorf("Caption for image %s does not exist", image.Filename)
		}(data.Images[index])
	}

	wg.Wait()
}

func (data *DataSet) appendCaptionsConcurrently() {
	startTime := time.Now()
	var wg sync.WaitGroup

	// Added using DataSet's locks
	data.captionLock.RLock()
	keys := make([]string, 0, len(data.TempCaption))
	for k := range data.TempCaption {
		keys = append(keys, k)
	}
	data.captionLock.RUnlock()

	// Iterate over the keys
	for _, k := range keys {
		// For each key, a new goroutine is launched to do the appending
		wg.Add(1)
		go data.appendCaption(&wg, k)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	elapsedTime := time.Since(startTime)
	captionLogPrinter.Infof("Finished appending captions concurrently in: %s", elapsedTime)
}

func (data *DataSet) appendCaption(waitGroup *sync.WaitGroup, key string) {
	defer waitGroup.Done()

	// Added using DataSet's locks
	data.captionLock.RLock()
	caption := data.TempCaption[key]
	data.captionLock.RUnlock()

	fileName, _ := strings.CutSuffix(caption.Filename, caption.Extension)
	tempCaptionLogPrinter := roggy.Printer(fmt.Sprintf("caption-handler: %v", caption.Filename))

	data.imagesLock.RLock()
	img, ok := data.Images[fileName]
	data.imagesLock.RUnlock()

	if !ok {
		tempCaptionLogPrinter.Errorf("Image file for caption %s does not exist", caption.Filename)
		return
	}

	data.captionLock.Lock()
	defer data.captionLock.Unlock()
	if img.Caption.Directory == caption.Directory {
		tempCaptionLogPrinter.Noticef("Caption file for image %s already exists", fileName)
		delete(data.TempCaption, fileName)
		return
	}

	if img.Caption.Directory != "" && img.Caption.Directory != caption.Directory {
		tempCaptionLogPrinter.Noticef("Caption file for image %s already exists but directories do not match", fileName)
		tempCaptionLogPrinter.Debugf("Image directory: %s", img.Directory)
		tempCaptionLogPrinter.Debugf("Old Caption directory: %s", img.Caption.Directory)
		tempCaptionLogPrinter.Debugf("New Caption directory: %s", caption.Directory)
		tempCaptionLogPrinter.Noticef("Assuming that we want the new caption file")
	}

	tempCaptionLogPrinter.Infof("Appending the caption file: %s to the image file: %s", caption.Filename, img.Filename)
	tempCaptionLogPrinter.Debugf("Caption directory: %s", caption.Directory)
	tempCaptionLogPrinter.Debugf("Image directory: %s", img.Directory)

	img.Caption = *caption
	data.imagesLock.Lock()
	data.Images[fileName] = img
	data.imagesLock.Unlock()

	delete(data.TempCaption, fileName)
}
