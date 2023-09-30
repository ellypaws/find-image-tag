package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const desiredSteps = 200

func addMultiple(steps int) tea.Cmd {
	if steps == 0 {
		steps = desiredSteps
	}
	return func() tea.Msg {
		return addMultipleMsg{current: 1, total: steps} // Start with the first step and desired total steps
	}
}

func directoryParse(filter int, directory string) tea.Cmd {
	// get all folders in directory (just one level deep)
	var dirEntries []os.DirEntry
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirEntries = append(dirEntries, entry)
		}
	}

	// count the number of directories
	numDirs := len(dirEntries)
	dirSliceFullPath := make([]string, numDirs)
	for i, dirEntry := range dirEntries {
		dirSliceFullPath[i] = filepath.Join(directory, dirEntry.Name())
	}

	msg := func() tea.Msg {
		return progressMsg{current: 1, total: numDirs, dirs: dirSliceFullPath}
	}

	return msg

}

func processDirectory(m *model, filter int, directory string, progress progressMsg) tea.Cmd {
	m.DataSet.WriteFiles(filter, directory)
	return func() tea.Msg {
		return progressMsg{
			current: progress.current + 1,
			total:   progress.total,
			dirs:    progress.dirs,
		}
	}
}

func addOne(num string) string {
	num = strings.Replace(num, ",", "", -1)
	newNum, _ := strconv.Atoi(num)
	newNum++
	newNumString := formatWithComma(newNum)
	return newNumString
}

func addPopulation(population string) tea.Cmd {
	newPopString := addOne(population)
	return func() tea.Msg {
		return popMsg(newPopString)
	}
}

func addCountImages(count string, row int) tea.Cmd {
	newCountString := addOne(count)
	updateNumReturn := updateNum{tableID: 0, row: row, column: 0, num: newCountString}
	return func() tea.Msg {
		return updateNumReturn
	}
}

func formatWithComma(i int) string {
	in := strconv.Itoa(i)
	var out strings.Builder
	l := len(in)
	for i, r := range in {
		_, _ = out.WriteRune(r)
		if (l-i-1)%3 == 0 && i < l-1 {
			_, _ = out.WriteRune(',')
		}
	}
	return out.String()
}
