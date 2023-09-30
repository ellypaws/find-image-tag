package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type CountFunction func(m model, tableID int, row int, column int) tea.Cmd
type CountRow []CountFunction
type Keys []CountRow

type EnterFunction func() tea.Cmd
type EnterActions []EnterFunction

// {CountFunction, CountFunction, CountFunction} --> countRow  ||
// {CountFunction, CountFunction, CountFunction} -- > countRow || --> keys
// {CountFunction, CountFunction, CountFunction} -- > countRow ||

func (m model) NewMenu() []Menu {
	stats, statsKeys, statsEnter := m.statsTable()
	captions, captionsKeys, captionsEnter := m.captionsTable()
	actions, actionsKeys, actionsEnter := m.actionsTable()

	return []Menu{
		{stats, statsKeys, statsEnter},
		{captions, captionsKeys, captionsEnter},
		{actions, actionsKeys, actionsEnter},
	}
}

func (m model) statsTable() (tbl table.Model, keys Keys, enter EnterActions) {
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

	values := Keys{
		{CountImagesWithCaptions},
		{CountCaptionDirectoryMatchImageDirectory},
		{CountImagesWithoutCaptions},
		{CountPending},
	}

	function := EnterActions{
		func() tea.Cmd { return nil },
		nil,
		nil,
		nil,
	}

	return table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(4),
	), values, function
}

func (m model) captionsTable() (tbl table.Model, keys Keys, enter EnterActions) {
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

	values := Keys{
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

	function := EnterActions{
		func() tea.Cmd { return nil },
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}

	return table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	), values, function
}

func (m model) actionsTable() (tbl table.Model, keys Keys, enter EnterActions) {
	columns := []table.Column{
		{Title: "Current", Width: 6},
		{Title: "New", Width: 6},
		{Title: "Total", Width: 6},
		{Title: "Actions", Width: 50},
	}

	rows := []table.Row{
		{"0", "0", "0", "Move captions to the image files"},
		{"0", "0", "0", "Hardlink captions to the image files"},
		{" ", " ", "0", "Merge new captions to existing captions"},
		{" ", " ", "0", "Append new tags to captions (dir)"},
		{" ", " ", " ", "Replace spaces with [_]"},
	}

	values := Keys{
		{CountImagesWithCaptionsNextToThem, CountOverwrites, CountImagesWithCaptions},
		{CountImagesWithCaptionsNextToThem, CountOverwrites, CountImagesWithCaptions},
		{nul, nul, CountCaptionsToMerge},
		{nul, nul, CountImagesWithCaptionsNextToThem},
		{nul, nul, nul},
	}

	function := EnterActions{
		func() tea.Cmd { return nil },
		nil,
		nil,
		nil,
		nil,
	}

	return table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	), values, function
}
