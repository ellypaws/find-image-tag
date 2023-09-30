package tui

import (
	"find-image-tag/tui/autocomplete"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	padding  = 2
	maxWidth = 80
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case addMultipleMsg:
		var numString string
		var modelToUpdate *table.Model
		var toUpdate *string

		// get the cell's address
		if m.menu[m.activeMenu].Focused() {
			modelToUpdate = &m.menu[m.activeMenu]
			toUpdate = &modelToUpdate.SelectedRow()[0]
		} else if m.table.Focused() {
			modelToUpdate = &m.table
			toUpdate = &modelToUpdate.SelectedRow()[3]
		}
		numString = *toUpdate
		numString = strings.Replace(numString, ",", "", -1)
		newNum, _ := strconv.Atoi(numString)
		newNum++
		newNumString := formatWithComma(newNum)

		// change the desired cell using the pointer
		*toUpdate = newNumString
		//m.menu[m.activeMenu].SelectedRow()[0] = newNumString
		//m.menu[m.activeMenu].UpdateViewport() // this is how we update after
		modelToUpdate.UpdateViewport()

		if msg.current%10 == 0 { // Only update progress for multiples of 5
			currentProgress := float64(msg.current) / float64(msg.total)
			cmd = m.progress.SetPercent(currentProgress)
		}

		if msg.current < msg.total {
			return m, tea.Batch(tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
				return addMultipleMsg{current: msg.current + 1, total: msg.total}
			}), cmd)
		} else {
			m.showProgress = false
		}
		return m, cmd
	case popMsg:
		s := string(msg)
		m.table.SelectedRow()[3] = s
		m.table.UpdateViewport() // this is how we update after

	// menu handlers
	case updateNum:
		m.menu[msg.tableID].Rows()[msg.row][msg.column] = msg.num
		m.menu[msg.tableID].UpdateViewport()
	//case countImagesWithCaptions:
	//	// set the first row's first column to the new count
	//	m.menu[0].Rows()[0][0] = string(msg)
	//	m.menu[0].UpdateViewport()
	//case countCaptionDirectoryMatchImageDirectory:
	//case countImagesWithoutCaptions:
	//case countPending:
	//case countFiles:
	//case countImages:
	//case countOverwrites:
	//case countCaptionsToMerge:
	//case countTotalCaptions:
	//case countImagesWithCaptionsNextToThem:
	//case offSet:
	//case moveString:

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case percentMsg:
		p := float64(msg)
		cmd = m.progress.SetPercent(p)
		return m, cmd

	case tea.KeyMsg:
		// Directory handler
		if m.showTextInput {
			if msg.Type == tea.KeyEnter {
				m.showTextInput = false
				m.table.Focus()
				return m, tea.Println(m.textInput.Value())
			}
			path := m.textInput.Value()
			if path == "" {
				path = m.textInput.Placeholder
			}

			basePath := filepath.Dir(path) // Gets the path up to the last backslash

			if _, err := os.Stat(basePath); err == nil {
				directories, err := autocomplete.GetDirectories(basePath)
				if err == nil {
					for i, dir := range directories {
						directories[i] = filepath.Join(basePath, dir)
					}
					m.textInput.SetSuggestions(directories)
				}
			} else {
				m.textInput.SetSuggestions(nil)
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":

			// chosen menu handler
			if m.menu[0].Focused() {
				return m, addCountImages(m.menu[0].Rows()[0][0], m.menu[0].Cursor())
			}
			if m.menu[1].Cursor() == 0 {
				m.showTextInput = true
				return m, nil
			}

			if m.table.Focused() && m.table.Cursor() == 1 {
				m.table.Blur()
				m.showTextInput = true
			}
			return m, tea.Batch(
				//tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
				addPopulation(m.table.SelectedRow()[3]),
			)
		case "a":
			if m.showTextInput {
				return m, nil
			}
			m.showProgress = true
			m.progress = progress.New(progress.WithDefaultGradient())
			runCmd := tea.Batch(addMultiple())
			return m, runCmd
		case "u":
			for menuID := range m.menu {
				for row := range m.menu[menuID].Rows() {
					cmd := m.keys[menuID][row]
					if cmd != nil {
						return m, cmd
					}
				}
			}
		}

		// menu movement handler
		for i, _ := range m.menu {
			// Based on focus state of tables, update the focused table
			if m.menu[i].Focused() {
				//m.menu[i], cmd = m.menu[i].Update(msg)
				if i < len(m.menu)-1 && msg.String() == "down" && m.menu[i].Cursor() == len(m.menu[i].Rows())-1 {
					m.menu[i].Blur()
					m.menu[i+1].Focus()
					m.menu[i+1].SetCursor(0)
					m.menu[i+1].SetStyles(focused)
					m.activeMenu = i + 1
					return m, nil
				}
				if i > 0 && msg.String() == "up" && m.menu[i].Cursor() == 0 {
					m.menu[i].Blur()
					m.menu[i-1].Focus()
					m.activeMenu = i - 1
					return m, nil
				}
			}
		}
	}

	// update tables
	var batch []tea.Cmd
	m.table, cmd = m.table.Update(msg)
	for i, _ := range m.menu {
		if m.menu[i].Focused() {
			m.menu[i], cmd = m.menu[i].Update(msg)
			batch = append(batch, cmd)
		}
	}
	return m, tea.Batch(append(batch, cmd)...)
}

type updateNum struct {
	num     string
	tableID int
	row     int
	column  int
}
type popMsg string
type percentMsg float64
type addMultipleMsg struct {
	current int
	total   int
}
