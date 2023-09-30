package autocomplete

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	textInput textinput.Model
}

func initialModel() model {
	ti := textinput.New()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	ti.Prompt = ""
	ti.Placeholder = currentDir + string(os.PathSeparator)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	ti.ShowSuggestions = true
	ti.SetSuggestions([]string{currentDir}) // Set initial suggestions to the current directory

	return model{textInput: ti}
}

func Init() textinput.Model {
	ti := textinput.New()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	ti.Prompt = ""
	ti.Placeholder = currentDir + string(os.PathSeparator)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	ti.ShowSuggestions = true
	ti.SetSuggestions([]string{currentDir}) // Set initial suggestions to the current directory

	return ti
}

func GetDirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}
	return directories, nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			fmt.Println(m.textInput.Value())
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		// Check for directory suggestions after key input
		path := m.textInput.Value()
		if path == "" {
			path = m.textInput.Placeholder
		}

		basePath := filepath.Dir(path) // Gets the path up to the last backslash

		if _, err := os.Stat(basePath); err == nil {
			directories, err := GetDirectories(basePath)
			if err == nil {
				for i, dir := range directories {
					directories[i] = filepath.Join(basePath, dir)
				}
				m.textInput.SetSuggestions(directories)
			}
		} else {
			m.textInput.SetSuggestions(nil)
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Choose a directory:\n\n%s\n\n%s\n",
		m.textInput.View(),
		"(tab to complete, ctrl+n/ctrl+p to cycle through suggestions, enter to select, esc to quit)",
	)
}
