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

func addCountImages(count string) tea.Cmd {
	newCountString := addOne(count)
	return func() tea.Msg {
		return countImagesWithCaptions(newCountString)
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
