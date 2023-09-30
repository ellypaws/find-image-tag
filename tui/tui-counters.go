package tui

import (
	"fmt"
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

func nul(m model, tableID int, row int, column int) tea.Cmd {
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: " "}
	}
}

func CountFiles(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountFiles())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountPending(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountPending())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountImages(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountImages())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountImagesWithCaptions(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountImagesWithCaptions())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountImagesWithCaptionsNextToThem(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountImagesWithCaptionsNextToThem())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountOverwrites(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountOverwrites(m.overwrite))
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountCaptionsToMerge(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountCaptionsToMerge())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountImagesWithoutCaptions(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountImagesWithoutCaptions())
	fmt.Println("CountImagesWithoutCaptions", newCount)
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountCaptionDirectoryMatchImageDirectory(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountCaptionDirectoryMatchImageDirectory())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}

func CountTotalCaptions(m model, tableID int, row int, column int) tea.Cmd {
	newCount := formatWithComma(m.DataSet.CountPending() + m.DataSet.CountImagesWithCaptions())
	return func() tea.Msg {
		return updateNum{tableID: tableID, row: row, column: column, num: newCount}
	}
}
