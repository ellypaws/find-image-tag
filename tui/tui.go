package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

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
