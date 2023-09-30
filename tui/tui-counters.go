package tui

import (
	"find-image-tag/entities"
	tea "github.com/charmbracelet/bubbletea"
)

type countImagesWithCaptions string
type countCaptionDirectoryMatchImageDirectory string
type countImagesWithoutCaptions string
type countPending string
type countFiles string
type countImages string
type countOverwrites string
type countCaptionsToMerge string
type countTotalCaptions string
type countImagesWithCaptionsNextToThem string
type offSet string
type moveString string

func (m model) nul(tableID int, row int, column int) tea.Cmd {
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: ""}
	}
}

func (m model) CountFiles(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountFiles())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountPending(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountPending())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountImages(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountImages())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountImagesWithCaptions(tableID int, row int, column int) tea.Cmd {
	m.DataSet = entities.InitDataSet()
	newCount := formatWithComma(m.DataSet.CountImagesWithCaptions())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountImagesWithCaptionsNextToThem(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountImagesWithCaptionsNextToThem())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountOverwrites(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountOverwrites(m.overwrite))
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountCaptionsToMerge(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountCaptionsToMerge())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountImagesWithoutCaptions(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountImagesWithoutCaptions())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountCaptionDirectoryMatchImageDirectory(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountCaptionDirectoryMatchImageDirectory())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func (m model) CountTotalCaptions(tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountPending() + m.DataSet.CountImagesWithCaptions())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

//package main
//
//import (
//	"os"
//	"path/filepath"
//)
//
//func (data *DataSet) countFiles() int {
//	return data.countPending() + data.countImages()
//}
//
//func (data *DataSet) countPending() int {
//	return len(data.TempCaption)
//}
//
//func (data *DataSet) countImages() int {
//	return len(data.Images)
//}
//
//func (data *DataSet) countImagesWithCaptions() int {
//	count := 0
//	for _, image := range data.Images {
//		if image.Caption.Filename != "" {
//			count++
//		}
//	}
//	return count
//}
//
//func (data *DataSet) countImagesWithCaptionsNextToThem() int {
//	count := 0
//	for name, image := range data.Images {
//		filePath := filepath.Join(image.Directory, name+".txt")
//		_, err := os.Stat(filePath)
//		if err == nil {
//			roggyPrinter.Debugf("File %s exists", filePath)
//			count++
//		}
//	}
//	return count
//}
//
//func (data *DataSet) countOverwrites(overwrite bool) int {
//	count := 0
//	for _, image := range data.Images {
//
//		// first check if image.Caption.Filename is empty
//		if image.Caption.Filename == "" {
//			continue
//		}
//
//		if !overwrite {
//			filePath := filepath.Join(image.Directory, image.Caption.Filename)
//			_, err := os.Stat(filePath)
//			if err == nil {
//				continue
//			}
//		}
//
//		// now check if caption directory is different from image directory
//
//		if image.Caption.Filename != "" && image.Caption.Directory != image.Directory {
//			count++
//		}
//	}
//	return count
//}
//
//// We're counting if the caption file exists literally next to
//// the image file. And then we check if the image and captions
//// are in different directories. That means we're appending
//// a new caption in place of the one that already exists alongside
//// the caption that already exists.
////
//// 📂 image
////
//// └── 📄 image.jpg
////
//// └── 📄 image.txt
////
//// 📂 new-caption
////
//// └── 📄 image.txt
//func (data *DataSet) countCaptionsToMerge() int {
//	count := 0
//	for _, image := range data.Images {
//
//		if image.Caption.Filename == "" {
//			continue
//		}
//
//		if image.Caption.Directory == image.Directory {
//			continue
//		}
//
//		filePath := filepath.Join(image.Directory, image.Caption.Filename)
//		_, err := os.Stat(filePath)
//		if err != nil {
//			// we can't merge if the caption file doesn't exist
//			continue
//		}
//
//		if image.Caption.Directory != image.Directory {
//			count++
//		}
//	}
//	return count
//}
//
//func (data *DataSet) countImagesWithoutCaptions() int {
//	count := 0
//	for _, image := range data.Images {
//		if image.Caption.Filename == "" {
//			count++
//		}
//	}
//	return count
//}
//
//func (data *DataSet) countCaptionDirectoryMatchImageDirectory() int {
//	count := 0
//	for _, image := range data.Images {
//		if image.Caption.Directory == image.Directory {
//			count++
//		}
//	}
//	return count
//}
