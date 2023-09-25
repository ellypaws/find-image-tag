package main

import (
	"fmt"
	"os"
	"regexp"
)

var imageDirectory string
var captionDirectory string

func main() {
	// Read all the filenames of image files in the image directory

	data := DataSet{
		Images: make(map[string]Image),
	}

	data.WriteFiles()
	fmt.Printf("Images: %v\n", data.Images)
	data.WriteFiles()
	fmt.Printf("Images: %v\n", data.Images)
	//data.CheckIfCaptionsExist()
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

	dirEntry, _ := os.ReadDir(directoryToRead)
	for _, entry := range dirEntry {
		t := entry.Type()
		fmt.Println(t.String())
		if entry.IsDir() {
			continue
		}
		currentEntry := entry.Name()
		extRegex, _ := regexp.Compile(regex)
		extension := extRegex.FindString(currentEntry)
		fileName := currentEntry[:len(currentEntry)-len(extension)]
		if !extRegex.MatchString(currentEntry) {
			continue
		}
		// check extension

		if extension == ".txt" {
			// check if image exists
			if _, ok := data.Images[fileName]; !ok {
				fmt.Println("Image file for caption", currentEntry, "does not exist")
			}
			if img, ok := data.Images[fileName]; ok {
				fmt.Println("Appending the caption file:", currentEntry, "to the image file:", fileName)
				img.Caption = Caption{Filename: currentEntry, Extension: extension}
				data.Images[fileName] = img
			}
			continue
		}

		data.Images[fileName] = Image{Filename: currentEntry, Caption: Caption{}}

	}
}

func (data *DataSet) WriteImages() {
	fmt.Print("Enter the path to the image directory: ")
	_, _ = fmt.Scanln(&imageDirectory)
	if imageDirectory == "" {
		imageDirectory = "."
	}

	dirEntry, _ := os.ReadDir(imageDirectory)
	for _, entry := range dirEntry {
		filename := entry.Name()
		// check if an image file
		extRegex, _ := regexp.Compile(`(?i)\.(jpe?g|png|gif|bmp)$`)
		if !extRegex.MatchString(filename) || entry.IsDir() {
			continue
		}
		data.Images[filename] = Image{Filename: filename}
	}
}

func (data *DataSet) WriteCaptions() {
	fmt.Print("Enter the path to the caption directory: ")
	_, _ = fmt.Scanln(&imageDirectory)
	if imageDirectory == "" {
		imageDirectory = "."
	}

	dirEntry, _ := os.ReadDir(imageDirectory)
	for _, entry := range dirEntry {
		filename := entry.Name()
		// check if an image file
		extRegex, _ := regexp.Compile(`(?i)\.(txt)$`)
		if !extRegex.MatchString(entry.Name()) || entry.IsDir() {
			continue
		}
		if img, ok := data.Images[filename]; ok {
			img.Caption = Caption{Filename: filename}
		}
	}
}
