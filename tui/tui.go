package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table        table.Model
	progress     progress.Model
	showProgress bool
}

func (m model) Init() tea.Cmd {
	m.progress = progress.New(progress.WithDefaultGradient())
	m.progress.Width = 50 // Setting initial width of progress bar
	return nil
}

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
			return m, tea.Batch(
				//tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
				addPopulation(m.table.SelectedRow()[3]),
			)
		case "a":
			m.showProgress = true
			m.progress = progress.New(progress.WithDefaultGradient())
			runCmd := tea.Batch(addMultiple())
			return m, runCmd
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

const desiredSteps = 200

func addMultiple() tea.Cmd {
	return func() tea.Msg {
		return addMultipleMsg{current: 1, total: desiredSteps} // Start with the first step and desired total steps
	}
}

func addPopulation(population string) tea.Cmd {
	// remove comma from string
	population = strings.Replace(population, ",", "", -1)
	newPop, _ := strconv.Atoi(population)
	newPop += 1
	newPopString := formatWithComma(newPop)
	return func() tea.Msg {
		return popMsg(newPopString)
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

func (m model) View() string {
	view := m.table.View() + "\n"
	if m.showProgress {
		view += m.progress.View() + "\n"
	}
	return baseStyle.Render(view) + "\n"
}

func Main() {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}

	rows := []table.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t, progress: progress.New(progress.WithDefaultGradient())}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
