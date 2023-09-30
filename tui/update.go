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
	AddBoth = iota
	AddCaption
	AddImage
)

const (
	padding  = 2
	maxWidth = 80
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var i directoryPrompt
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case progressMsg:
		if msg.current%10 == 0 { // Only update progress for multiples of 5
			currentProgress := float64(msg.current) / float64(msg.total)
			cmd = m.progress.SetPercent(currentProgress)
		}

		if msg.current < msg.total {
			processDirectory(&m, AddBoth, msg.dirs[msg.current-1], msg)
			return m, tea.Batch(tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
				return progressMsg{current: msg.current + 1, total: msg.total, dirs: msg.dirs}
			}), cmd)
		} else {
			m.showProgress = false
		}
		return m, cmd
	case addMultipleMsg:
		var numString string
		var modelToUpdate *table.Model
		var toUpdate *string

		// get the cell's address
		for menuID, currentMenu := range m.menus {
			if currentMenu.Menu.Focused() {
				modelToUpdate = &m.menus[menuID].Menu
				toUpdate = &modelToUpdate.SelectedRow()[0]
			}
		}

		if m.table.Focused() {
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
	case msgToPrint:
		return m, tea.Printf(string(msg))
	case directoryPrompt:
		i = msg
		m.showTextInput = true
	case []Menu:
		m.menus = msg
		for menuID, _ := range m.menus {
			m.menus[menuID].Menu.UpdateViewport()
		}
	case updateNum: // menu handlers
		m.menus[msg.tableID].Menu.Rows()[msg.row][msg.column] = msg.num
		m.menus[msg.tableID].Menu.UpdateViewport()
	case startCount:
		if msg {
			for menuID, currentMenu := range m.menus {
				for row := range currentMenu.Menu.Rows() {
					for column := range currentMenu.Menu.Rows()[row] {
						cmdFunc := currentMenu.UpdateFunc[row][column]
						if cmdFunc != nil {
							cmd := cmdFunc(m, menuID, row, column)
							return m, cmd
						}
					}
				}
			}
		}
	case progress.FrameMsg: // FrameMsg is sent when the progress bar wants to animate itself
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case percentMsg:
		p := float64(msg)
		cmd = m.progress.SetPercent(p)
		return m, cmd
	case tea.KeyMsg:
		if m.showTextInput { // Directory handler
			if msg.Type == tea.KeyEnter {
				m.showTextInput = false
				m.menus[0].Menu.Focus()
				m.DataSet.WriteFiles(int(i), m.textInput.Value())
				return m, Refresh()
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
			for _, currentMenu := range m.menus {
				if currentMenu.Menu.Focused() {
					if ok := currentMenu.EnterFunc[currentMenu.Menu.Cursor()]; ok != nil {
						cmd = ok(m, cmd)
						return m, cmd
					}
				}
			}
			return m, Refresh()
		case "a":
			if m.showTextInput {
				return m, nil
			}
			m.showProgress = true
			m.progress = progress.New(progress.WithDefaultGradient())
			runCmd := tea.Batch(addMultiple(desiredSteps))
			return m, runCmd
		case "u":
			return m, Refresh()
		}

		for menuID, currentMenu := range m.menus { // up/down handlers
			if currentMenu.Menu.Focused() {
				//m.menu[i], cmd = m.menu[i].Update(msg)
				if menuID < len(m.menus)-1 && msg.String() == "down" && currentMenu.Menu.Cursor() == len(currentMenu.Menu.Rows())-1 {
					m.menus[menuID].Menu.Blur()
					m.menus[menuID].Menu.SetStyles(unfocused)
					m.menus[menuID+1].Menu.Focus()
					m.menus[menuID+1].Menu.SetStyles(focused)
					m.menus[menuID+1].Menu.SetCursor(0)
					return m, nil
				}
				if menuID > 0 && msg.String() == "up" && currentMenu.Menu.Cursor() == 0 {
					m.menus[menuID].Menu.Blur()
					m.menus[menuID].Menu.SetStyles(unfocused)
					m.menus[menuID-1].Menu.Focus()
					m.menus[menuID-1].Menu.SetStyles(focused)
					m.menus[menuID-1].Menu.SetCursor(len(m.menus[menuID-1].Menu.Rows()) - 1)
					return m, nil
				}
			}
		}
	}

	// update tables
	var batch []tea.Cmd
	m.table, cmd = m.table.Update(msg)

	for menuID, currentMenu := range m.menus {
		m.menus[menuID].Menu, cmd = currentMenu.Menu.Update(msg)
		batch = append(batch, cmd)
	}
	return m, tea.Batch(append(batch, cmd)...)
}

func Refresh() tea.Cmd {
	return func() tea.Msg {
		return startCount(true)
	}
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
type progressMsg struct {
	current int
	total   int
	dirs    []string
}
type msgToPrint string
type startCount bool

type directoryPrompt int

const (
	CheckExist = iota
	Print
	Reset
	WriteJSON
	Append
	CheckMissing
	Quit
)

const (
	MoveCaptions = iota
	Hardlink
	Merge
	AddTags
	Underscores
)

type Actions struct {
	Menu int
	Do   int
}
