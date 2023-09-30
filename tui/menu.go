package tui

import "github.com/charmbracelet/bubbles/table"

func NewMenu() []table.Model {
	var menu []table.Model
	menu = append(menu, statsTable())
	menu = append(menu, captionsTable())
	menu = append(menu, actionsTable())
	return menu
}

func statsTable() table.Model {
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

	return table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(4),
	)
}

func captionsTable() table.Model {
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

	return table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)
}

func actionsTable() table.Model {
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
		{"0", "0", "0", "Append new tags to captions (dir)"},
		{" ", " ", " ", "Replace spaces with [_]"},
	}

	return table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)
}
