package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (data *DataSet) CaptionsToImages(move bool) {
	for _, image := range data.Images {
		if image.Caption.Filename == "" {
			roggyPrinter.Noticef("Image %s does not have a caption", image.Filename, "Skipping...")
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
				roggyPrinter.Errorf("Error moving file: %v", err)
			} else {
				roggyPrinter.Noticef("File moved successfully from %s to %s", from, to)
			}
		} else {
			// Create a hardlink of the file.
			err := os.Link(from, to)
			if err != nil {
				roggyPrinter.Errorf("Error linking file: %v", err)
			} else {
				roggyPrinter.Infof("File linked successfully from %s to %s", from, to)
			}
		}
	}
}

func (data *DataSet) replaceSpaces() {
	re := regexp.MustCompile(`(\w)\s+(\w)`)

	for _, image := range data.Images {
		captionFile := filepath.Join(image.Caption.Directory, image.Caption.Filename)

		file, err := os.Open(captionFile)
		if err != nil {
			roggyPrinter.Errorf("Error opening file: %v", err)
			continue
		}

		reader := bufio.NewReader(file)
		content, _ := reader.ReadString('\n')
		file.Close()

		newContent := re.ReplaceAllString(content, "${1}_${2}")

		err = os.Remove(captionFile)
		if err != nil {
			roggyPrinter.Errorf("Error deleting file: %v", err)
			continue
		}

		file, err = os.Create(captionFile)
		if err != nil {
			roggyPrinter.Errorf("Error creating file: %v", err)
			continue
		}

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(newContent + "\n")
		writer.Flush()
		file.Close()

		if err != nil {
			roggyPrinter.Errorf("Error writing file: %v", err)
		} else {
			roggyPrinter.Infof("Replaced spaces with underscores for file: %s", captionFile)
		}
	}
}

func (data *DataSet) CheckIfCaptionsExist() {
	for _, image := range data.Images {
		if image.Caption.Filename != "" {
			continue
		}
		roggyPrinter.Errorf("Caption for image %s does not exist", image.Filename)
	}
}

func (data *DataSet) checkForMissingImages() {
	for _, caption := range data.TempCaption {
		fileName, _ := strings.CutSuffix(caption.Filename, caption.Extension)
		if _, ok := data.Images[fileName]; !ok {
			roggyPrinter.Errorf("Image file for caption %s does not exist", caption.Filename)
		}
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
			if _, ok := data.TempCaption[fileName]; !ok {
				if data.TempCaption == nil {
					data.InitCaption()
				}
				newCaption := Caption{Filename: currentEntry, Extension: extension, Directory: directory}
				data.TempCaption[fileName] = &newCaption
				roggyPrinter.Infof("Added file: %s to the temporary caption dataset", currentEntry)
				roggyPrinter.Debugf("Directory: %s", directory)
			} else {
				roggyPrinter.Noticef("Caption file in temporary caption dataset for image %s already exists", fileName)
			}
			return nil
		}

		if _, ok := data.Images[fileName]; !ok {
			newImg := Image{Filename: currentEntry, Extension: extension, Directory: directory, Caption: Caption{}}
			data.Images[fileName] = &newImg
			roggyPrinter.Infof("Added file: %s to the dataset", currentEntry)
			roggyPrinter.Debugf("Directory: %s", directory)
		}

		return nil
	})
	if err != nil {
		return
	}

	data.appendCaptions()
}

func (data *DataSet) appendCaptions() {
	for _, caption := range data.TempCaption {
		fileName, _ := strings.CutSuffix(caption.Filename, caption.Extension)
		if img, ok := data.Images[fileName]; ok {
			roggyPrinter.Infof("Appending the caption file: %s to the image file: %s", caption.Filename, fileName)
			roggyPrinter.Debugf("Directory: %s", caption.Directory)
			img.Caption = *caption
			data.Images[fileName] = img

			// now remove from the temp caption dataset
			delete(data.TempCaption, fileName)
		} else {
			roggyPrinter.Noticef("Image file for caption %s does not exist", caption.Filename)
		}
	}
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
