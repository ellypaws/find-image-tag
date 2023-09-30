package entities

import (
	"github.com/nokusukun/roggy"
	"os"
	"path/filepath"
)

var roggyPrinter = roggy.Printer("main-service")

func (data *DataSet) CountFiles() int {
	return data.CountPending() + data.CountImages()
}

func (data *DataSet) CountPending() int {
	return len(data.TempCaption)
}

func (data *DataSet) CountImages() int {
	return len(data.Images)
}

func (data *DataSet) CountImagesWithCaptions() int {
	count := 0
	for _, image := range data.Images {
		if image.Caption.Filename != "" {
			count++
		}
	}
	return count
}

func (data *DataSet) CountImagesWithCaptionsNextToThem() int {
	count := 0
	for name, image := range data.Images {
		filePath := filepath.Join(image.Directory, name+".txt")
		_, err := os.Stat(filePath)
		if err == nil {
			roggyPrinter.Debugf("File %s exists", filePath)
			count++
		}
	}
	return count
}

func (data *DataSet) CountOverwrites(overwrite bool) int {
	count := 0
	for _, image := range data.Images {

		// first check if image.Caption.Filename is empty
		if image.Caption.Filename == "" {
			continue
		}

		if !overwrite {
			filePath := filepath.Join(image.Directory, image.Caption.Filename)
			_, err := os.Stat(filePath)
			if err == nil {
				continue
			}
		}

		// now check if caption directory is different from image directory

		if image.Caption.Filename != "" && image.Caption.Directory != image.Directory {
			count++
		}
	}
	return count
}

// CountCaptionsToMerge We're counting if the caption file exists literally next to
// the image file. And then we check if the image and captions
// are in different directories. That means we're appending
// a new caption in place of the one that already exists alongside
// the caption that already exists.
//
// 📂 image
//
// └── 📄 image.jpg
//
// └── 📄 image.txt
//
// 📂 new-caption
//
// └── 📄 image.txt
func (data *DataSet) CountCaptionsToMerge() int {
	count := 0
	for _, image := range data.Images {

		if image.Caption.Filename == "" {
			continue
		}

		if image.Caption.Directory == image.Directory {
			continue
		}

		filePath := filepath.Join(image.Directory, image.Caption.Filename)
		_, err := os.Stat(filePath)
		if err != nil {
			// we can't merge if the caption file doesn't exist
			continue
		}

		if image.Caption.Directory != image.Directory {
			count++
		}
	}
	return count
}

func (data *DataSet) CountImagesWithoutCaptions() int {
	count := 0
	for _, image := range data.Images {
		if image.Caption.Filename == "" {
			count++
		}
	}
	return count
}

func (data *DataSet) CountCaptionDirectoryMatchImageDirectory() int {
	count := 0
	for _, image := range data.Images {
		if image.Caption.Directory == image.Directory {
			count++
		}
	}
	return count
}
