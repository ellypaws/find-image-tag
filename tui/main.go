package tui

import (
	"find-image-tag/entities"
	"find-image-tag/tui/autocomplete"
	"find-image-tag/tui/sender"
	textinputs "find-image-tag/tui/text-inputs"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"os"
	"time"
)

type model struct {
	DataSet              *entities.DataSet
	menus                []Menu
	overwrite            bool
	progress             progress.Model
	showProgress         bool
	textInput            textinput.Model
	showTextInput        bool
	multiTextInput       textinputs.Model
	showMultiInput       bool
	sender               sender.Model
	senderActiveDuration time.Duration
}

type Menu struct {
	Menu       table.Model
	UpdateFunc Keys
	EnterFunc  EnterActions
}

func Main() {
	focused = table.DefaultStyles()
	focused.Header = focused.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	focused.Selected = focused.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	unfocused = table.DefaultStyles()
	unfocused.Header = unfocused.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	unfocused.Selected = unfocused.Selected.
		Foreground(lipgloss.Color("#bbbbbb"))

	m := model{
		progress:       progress.New(progress.WithDefaultGradient()),
		textInput:      autocomplete.Init(),
		menus:          model{}.NewMenu(),
		DataSet:        entities.InitDataSet(),
		multiTextInput: textinputs.InitialModel(),
		sender:         sender.NewModel(),
	}

	// set styles for each menu
	for menuID, currentMenu := range m.menus {
		if currentMenu.Menu.Focused() {
			m.menus[menuID].Menu.SetStyles(focused)
		} else {
			m.menus[menuID].Menu.SetStyles(unfocused)
		}
	}

	p := tea.NewProgram(m)

	go func() {
		for i := 0; i < 3; i++ {
			m.sender.Active = true
			pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec
			time.Sleep(pause)

			// Send the Bubble Tea program a message from outside the
			// tea.Program. This will block until it is ready to receive
			// messages.
			p.Send(sender.ResultMsg{Food: sender.RandomFood(), Duration: pause})
		}
		m.sender.Active = false
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
