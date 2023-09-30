package tui

import (
	"find-image-tag/entities"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	m.progress = progress.New(progress.WithDefaultGradient())
	m.progress.Width = 50 // Setting initial width of progress bar
	m.textInput = textinput.New()
	m.DataSet = entities.InitDataSet()
	m.menus = m.NewMenu()
	return nil
}
