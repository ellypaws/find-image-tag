package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type CountFunction func(tableID int, row int, column int) tea.Cmd
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

	//menu := []table.Model{
	//	stats,
	//	captions,
	//	actions,
	//}
	//
	//keys := []Keys{
	//	statsKeys,
	//	captionsKeys,
	//	actionsKeys,
	//}
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
		{m.CountImagesWithCaptions},
		{m.CountCaptionDirectoryMatchImageDirectory},
		{m.CountImagesWithoutCaptions},
		{m.CountPending},
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

	//values := CountRow{
	//	tempM.CountFiles,
	//	tempM.CountTotalCaptions,
	//	tempM.CountImages,
	//	tempM.CountImages,
	//	tempM.nul,
	//}

	values := Keys{
		{m.CountFiles},
		{m.CountTotalCaptions},
		{m.CountImages},
		{m.CountImages},
		{m.nul},
		{m.nul},
		{m.nul},
		{m.CountPending},
		{m.nul},
		{m.nul},
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

	//values := CountRow{
	//	tempM.CountFiles,
	//	tempM.CountImages,
	//}

	values := Keys{
		{m.CountImagesWithCaptionsNextToThem, m.CountOverwrites, m.CountImagesWithCaptions},
		{m.CountImagesWithCaptionsNextToThem, m.CountOverwrites, m.CountImagesWithCaptions},
		{m.nul, m.nul, m.CountCaptionsToMerge},
		{m.nul, m.nul, m.CountImagesWithCaptionsNextToThem},
		{m.nul, m.nul, m.nul},
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
