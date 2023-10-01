package tui

import (
	"find-image-tag/entities"
	"find-image-tag/tui/sender"
	textinputs "find-image-tag/tui/text-inputs"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	m.progress = progress.New(progress.WithDefaultGradient())
	m.progress.Width = 50 // Setting initial width of progress bar
	m.textInput = textinput.New()
	m.DataSet = entities.InitDataSet()
	m.sender = sender.NewModel()
	m.multiTextInput = textinputs.InitialModel()
	m.menus = m.NewMenu()
	m.sender.Active = true
	return nil
}
