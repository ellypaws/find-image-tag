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

var captionLogPrinter = roggy.Printer("caption-handler")

func (data *DataSet) CaptionsToImages(move bool) {
	for _, image := range data.Images {
		if image.Caption.Filename == "" {
			captionLogPrinter.Noticef("Image %s does not have a caption", image.Filename, "Skipping...")
			continue
		}

		// Handle relative directories, get current directory with os.Getwd
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

		if move {
			// Move the file. On many file systems, this is a simple rename operation.
			err := os.Rename(from, to)
			if err != nil {
				captionLogPrinter.Errorf("Error moving file: %v", err)
			} else {
				captionLogPrinter.Noticef("File moved successfully from %s to %s", from, to)
			}
		} else {
			// Create a hardlink of the file.
			err := os.Link(from, to)
			if err != nil {
				captionLogPrinter.Errorf("Error linking file: %v", err)
			} else {
				captionLogPrinter.Infof("File linked successfully from %s to %s", from, to)
			}
		}
	}
}

func (data *DataSet) replaceSpaces() {
	re := regexp.MustCompile(`(\w)\s+(\w)`)
	wg := &sync.WaitGroup{}

	for _, image := range data.Images {
		i := image

		wg.Add(1)
		go func() {
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
			file.Close()

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
			writer.Flush()
			file.Close()

			if err != nil {
				captionLogPrinter.Errorf("Error writing file: %v", err)
			} else {
				captionLogPrinter.Infof("Replaced spaces with underscores for file: %s", captionFile)
			}
		}()
	}
	wg.Wait()
}

func (data *DataSet) CheckIfCaptionsExist() {
	for _, image := range data.Images {
		if image.Caption.Filename != "" {
			roggyPrinter.Noticef("Caption `%s` has a matching image file `%s`", image.Caption.Filename, image.Filename)
			roggyPrinter.Debugf("Caption directory: %s", image.Caption.Directory)
			roggyPrinter.Debugf("Image directory: %s", image.Directory)
			continue
		}
		roggyPrinter.Errorf("Caption for image %s does not exist", image.Filename)
	}
}

func (data *DataSet) checkForMissingImages() {
	wg := &sync.WaitGroup{}

	for ix := range data.TempCaption {
		wg.Add(1)
		go func(c *Caption) {
			defer wg.Done()
			tempCaptionLogPrinter := roggy.Printer(fmt.Sprintf("caption-handler: %v", c.Filename))
			fileName, _ := strings.CutSuffix(c.Filename, c.Extension)
			if _, ok := data.Images[fileName]; !ok {
				tempCaptionLogPrinter.Errorf("Image file for caption %s does not exist", c.Filename)
			}
		}(data.TempCaption[ix])
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

func (data *DataSet) countFiles() int {
	return data.countPending() + data.countImages()
}

func (data *DataSet) countPending() int {
	return len(data.TempCaption)
}

func (data *DataSet) countImages() int {
	return len(data.Images)
}

func (data *DataSet) countImagesWithCaptions() int {
	count := 0
	for _, image := range data.Images {
		if image.Caption.Filename != "" {
			count++
		}
	}
	return count
}

func (data *DataSet) countImagesWithoutCaptions() int {
	count := 0
	for _, image := range data.Images {
		if image.Caption.Filename == "" {
			count++
		}
	}
	return count
}

func (data *DataSet) countCaptionDirectoryMatchImageDirectory() int {
	count := 0
	for _, image := range data.Images {
		if image.Caption.Directory == image.Directory {
			count++
		}
	}
	return count
}
