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

func (data *DataSet) CaptionsToImages(action int, overwrite bool) {
	for _, image := range data.Images {
		if image.Caption.Filename == "" {
			captionLogPrinter.Noticef("NO CAPTION: Image %s does not have a caption. Skipping...", image.Filename)
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

		if from == to {
			captionLogPrinter.Noticef("SAME DIR: Directory of caption is already the same as the image. Skipping...")
			continue
		}

		switch action {
		case move:
			// Check if the file already exists
			_, err := os.Stat(to)
			if err == nil {
				if !overwrite {
					captionLogPrinter.Noticef("EXIST: File already exists. Skipping...")
					continue
				} else {
					captionLogPrinter.Noticef("OVERWRITE: File already exists. Overwriting...")
				}
			}

			// Move the file. On many file systems, this is a simple rename operation.
			err = os.Rename(from, to)
			if err != nil {
				captionLogPrinter.Errorf("Error moving file: %v", err)
			} else {
				captionLogPrinter.Noticef("File moved successfully from %s to %s", from, to)
			}
		case hardlink:
			// Check if the file already exists
			if !overwrite {
				_, err := os.Stat(to)
				if err == nil {
					captionLogPrinter.Noticef("EXIST: File already exists. Skipping...")
					continue
				}
			}

			if overwrite {
				_, err := os.Stat(to)
				if err == nil {
					captionLogPrinter.Noticef("EXIST: File already exists. Overwriting...")
					err = os.Remove(to)
					if err != nil {
						captionLogPrinter.Errorf("Error removing file: %v", err)
						continue
					}
				}
			}

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
		}
		// Update the image.Caption.Directory to the new directory 'to'
		image.Caption.Directory = filepath.Dir(to)
	}
}

func appendNewTags() {
	scanner := bufio.NewScanner(os.Stdin)

	// Prompt for a directory
	captionLogPrinter.Infof("Please enter directory path: ")
	scanner.Scan()
	directory := scanner.Text()

	// Prompt for new tags
	captionLogPrinter.Infof("Please enter new tags to add (separated by comma): ")
	scanner.Scan()
	newTagsInput := strings.TrimSpace(scanner.Text())

	// Split newTagsInput string to slice of tags
	newTags := strings.Split(newTagsInput, ",")
	for i, tag := range newTags {
		newTags[i] = strings.TrimSpace(strings.ToLower(tag))
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		captionLogPrinter.Errorf("Failed to read directory: %v", err)
		return
	}

	changeAll := false

	// Process each file in the directory
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".txt" {
			continue
		}

		// Read the existing tags from the file
		bs, err := os.ReadFile(filepath.Join(directory, f.Name()))
		if err != nil {
			captionLogPrinter.Errorf("Failed to read file: %v", err)
			continue
		}

		// Process the content to lower case
		content := strings.ToLower(string(bs))
		// Split the content, trim spaces, remove empty tags
		existingTags := strings.Split(content, ",")
		for i, tag := range existingTags {
			existingTags[i] = strings.TrimSpace(tag)
		}

		// Merge new tags and existing tags eliminating duplicates
		mergedTags := append(newTags)
		for _, existingTag := range existingTags {
			if existingTag == "" {
				continue
			}

			existsInNewTags := false
			for _, newTag := range newTags {
				if newTag == existingTag {
					existsInNewTags = true
					break
				}
			}

			if !existsInNewTags {
				mergedTags = append(mergedTags, existingTag)
			}
		}

		// Join the tags with commas
		newContent := strings.Join(mergedTags, ", ") + "\n"

		if !changeAll {
			// Preview changes and prompt for confirmation
			captionLogPrinter.Infof("New content for file %s will be:\n%s", f.Name(), newContent)
			captionLogPrinter.Infof("Apply changes? (y/n/all): ")
			scanner.Scan()
			userInput := scanner.Text()
			if userInput == "n" {
				captionLogPrinter.Infof("Changes not applied. Moving to next file.")
				continue
			} else if userInput == "all" {
				changeAll = true
			}
		}

		// Write the updated content to the file
		err = os.WriteFile(filepath.Join(directory, f.Name()), []byte(newContent), 0644)
		if err != nil {
			captionLogPrinter.Errorf("Failed to write to file: %v", err)
			continue
		}

		// Create a preview snippet of the new text
		end := len(newContent)
		desiredLength := len(newTagsInput) + 50
		if len(newContent) > desiredLength {
			end = desiredLength
		}
		previewSnippet := newContent[:end]

		captionLogPrinter.Infof("Changes applied successfully to %s", f.Name())
		captionLogPrinter.Debugf("\"%s...\"", previewSnippet)
	}
}

func (data *DataSet) replaceSpaces() {
	re := regexp.MustCompile(`(\w) ([(]?\w)`)
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
