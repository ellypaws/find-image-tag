package tui

import (
	"find-image-tag/tui/sender"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Menu struct {
	Menu          table.Model
	UpdateFunc    []CountRow
	EnterFunction []tea.Msg
}

type CountFunction func(m model, tableID int, row int, column int) tea.Cmd
type CountRow []CountFunction

// {CountFunction, CountFunction, CountFunction} --> countRow  ||
// {CountFunction, CountFunction, CountFunction} -- > countRow || --> keys
// {CountFunction, CountFunction, CountFunction} -- > countRow ||

func (m model) NewMenu() []Menu {

	return []Menu{
		m.statsTable(),
		m.captionsTable(),
		m.actionsTable(),
	}
}

func (m model) statsTable() Menu {
	columns := []table.Column{
		{Title: "#", Width: 6},
		{Title: "Stats", Width: 50},
	}

	rows := []table.Row{
		{"0", "Images with captions"},
		{"0", "Images with captions that match directories"},
		{"0", "Missing captions"},
		{"0", "Pending text files"},
	}

	values := []CountRow{
		{CountImagesWithCaptions},
		{CountCaptionDirectoryMatchImageDirectory},
		{CountImagesWithoutCaptions},
		{CountPending},
	}

	function := []tea.Msg{
		sender.ResultMsg{Food: "test"},
		startCount(true),
		startCount(true),
		startCount(true),
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(4),
	)

	return Menu{
		Menu:          t,
		UpdateFunc:    values,
		EnterFunction: function,
	}
}

func (m model) captionsTable() Menu {
	columns := []table.Column{
		{Title: "#", Width: 6},
		{Title: "Captions", Width: 50},
	}

	rows := []table.Row{
		{"0", "Add files to the dataset"},
		{"0", "Add captions to the dataset"},
		{"0", "Add images to the dataset"},
		{"0", "Check if each image has a caption"},
		{" ", "Print the dataset as JSON"},
		{" ", "Reset the dataset"},
		{" ", "Write the dataset as a JSON file"},
		{"0", "Append text files to matching images"},
		{" ", "Check for captions without matching images"},
		{" ", "Quit"},
	}

	values := []CountRow{
		{CountFiles},
		{CountTotalCaptions},
		{CountImages},
		{CountImages},
		{nul},
		{nul},
		{nul},
		{CountPending},
		{nul},
		{nul},
	}

	function := []tea.Msg{
		directoryPrompt(AddBoth),
		directoryPrompt(AddCaption),
		directoryPrompt(AddImage),
		Actions{CaptionsMenu, CheckExist},
		Actions{CaptionsMenu, Print},
		Actions{CaptionsMenu, Reset},
		Actions{CaptionsMenu, WriteJSON},
		Actions{CaptionsMenu, Append},
		Actions{CaptionsMenu, CheckMissing},
		Actions{CaptionsMenu, Quit},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	return Menu{
		Menu:          t,
		UpdateFunc:    values,
		EnterFunction: function,
	}

}

func (m model) actionsTable() Menu {
	columns := []table.Column{
		{Title: "Current", Width: 8},
		{Title: "New", Width: 6},
		{Title: "Total", Width: 8},
		{Title: "Actions", Width: 40},
	}

	rows := []table.Row{
		{"0", "0", "0", "Move captions to the image files"},
		{"0", "0", "0", "Hardlink captions to the image files"},
		{" ", " ", "0", "Merge new captions to existing captions"},
		{" ", " ", "0", "Append new tags to captions (dir)"},
		{" ", " ", " ", "Replace spaces with [_]"},
	}

	values := []CountRow{
		{CountImagesWithCaptionsNextToThem, CountOverwrites, CountImagesWithCaptions},
		{CountImagesWithCaptionsNextToThem, CountOverwrites, CountImagesWithCaptions},
		{nul, nul, CountCaptionsToMerge},
		{nul, nul, CountImagesWithCaptionsNextToThem},
		{nul, nul, nul},
	}

	function := []tea.Msg{
		Actions{ActionsMenu, MoveCaptions},
		Actions{ActionsMenu, Hardlink},
		Actions{ActionsMenu, Merge},
		Actions{ActionsMenu, AddTags},
		Actions{ActionsMenu, Underscores},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)

	return Menu{
		Menu:          t,
		UpdateFunc:    values,
		EnterFunction: function,
	}

}
