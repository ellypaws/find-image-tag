package main

import "os"

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

func (data *DataSet) countOverwrites() int {
	count := 0
	for _, image := range data.Images {
		_, err := os.Stat(image.Filename)
		if image.Caption.Filename != "" && err == nil {
			count++
		}
	}
	return count
}

func (data *DataSet) countExistingCaptions() int {
	// Use os.Stat
	count := 0
	for key, _ := range data.Images {
		if _, err := os.Stat(key + ".txt"); err == nil {
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
