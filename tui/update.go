package tui

import (
	"find-image-tag/tui/autocomplete"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case addMultipleMsg:
		population := m.table.SelectedRow()[3]
		population = strings.Replace(population, ",", "", -1)
		newPop, _ := strconv.Atoi(population)
		newPop++
		m.table.SelectedRow()[3] = formatWithComma(newPop)
		m.table.UpdateViewport()

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

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	//case resetMsg:
	//	m.progress = progress.New(progress.WithDefaultGradient())

	case percentMsg:
		p := float64(msg)
		cmd = m.progress.SetPercent(p)
		return m, cmd

	case tea.KeyMsg:
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
			if m.showTextInput {
				m.showTextInput = false
				m.table.Focus()
				tea.Println(m.textInput.Value())
				return m, nil
			}
			if m.table.Cursor() == 1 {
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
		}

		// Directory handler
		if m.showTextInput {
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
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

type popMsg string
type percentMsg float64
type addMultipleMsg struct {
	current int
	total   int
}
