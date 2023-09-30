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
		{Title: "#", Width: 4},
		{Title: "", Width: 20},
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
		table.WithHeight(7),
	)
}

func captionsTable() table.Model {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "", Width: 20},
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
		table.WithHeight(7),
	)
}

func actionsTable() table.Model {
	columns := []table.Column{
		{Title: "Current", Width: 4},
		{Title: "New", Width: 4},
		{Title: "Total", Width: 4},
		{Title: "", Width: 20},
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

//toPrint := []string{
//"1::" + roggy.Rainbowize("---") + " Stats " + roggy.Rainbowize("---"),
//"2::",
//stemp.Inline("2::{0:w=30,j=r} | Images with captions", countImagesWithCaptions),
//stemp.Inline("2::{0:w=30,j=r} | Images with captions that match directories", countCaptionDirectoryMatchImageDirectory),
//stemp.Inline("2::{0:w=30,j=r} | Missing captions", countImagesWithoutCaptions),
//stemp.Inline("2::{0:w=30,j=r} | Pending text files", countPending),
//"1::" + roggy.Rainbowize("---") + " Image Captioning " + roggy.Rainbowize("---"),
//"2::",
//stemp.Inline("2::{0:w=30,j=r} | [+] Add files to the dataset", countFiles),
//stemp.Inline("2::{0:w=30,j=r} | [+c] Add captions to the dataset", countTotalCaptions),
//stemp.Inline("2::{0:w=30,j=r} | [+i] Add images to the dataset", countImages),
//stemp.Inline("2::{0:w=30,j=r} | [C]heck if each image has a caption", countImages),
//stemp.Inline("2::{0:w=30,j=r} | [P]rint the dataset as JSON", nul),
//stemp.Inline("2::{0:w=30,j=r} | [R]eset the dataset", nul),
//stemp.Inline("2::{0:w=30,j=r} | [W]rite the dataset as a JSON file", nul),
//stemp.Inline("2::{0:w=30,j=r} | Append [t]ext files to matching images", countPending),
//stemp.Inline("2::{0:w=30,j=r} | Check for captions without matching [i]mages", nul),
//stemp.Inline("2::{0:w=30,j=r} | [Q]uit", nul),
//"1::" + roggy.Rainbowize("---") + " Actions " + roggy.Rainbowize("---"),
//"2::",
//stemp.Inline("2::{0:w=30,j=r} | {1} | {2:w=10,j=r}", ow, overwrite, overwriteString),
//stemp.Inline("2::"+moveString+" | [Move] {3}", countImagesWithCaptionsNextToThem, countOverwrites, countImagesWithCaptions, moveOverwriteString),
//stemp.Inline("2::"+moveString+" | [Hardlink] {3}", countImagesWithCaptionsNextToThem, countOverwrites, countImagesWithCaptions, moveOverwriteString),
//stemp.Inline("2::{0:w=30,j=r} | [Merge] new captions to existing captions", countCaptionsToMerge),
//stemp.Inline("2::{0:w=30,j=r} | [Append] new tags to captions (dir)", countImagesWithCaptionsNextToThem),
//stemp.Inline("2::{0:w=30,j=r} | Replace spaces with [_]", nul),
//}
