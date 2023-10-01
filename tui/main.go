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
	table                table.Model
	menus                []Menu
	overwrite            bool
	progress             progress.Model
	showProgress         bool
	textInput            textinput.Model
	showTextInput        bool
	multiTextInput       textinputs.Model
	showMultiInput       bool
	sender               sender.Model
	senderActiveDuration *time.Duration
}

type Menu struct {
	Menu       table.Model
	UpdateFunc Keys
	EnterFunc  EnterActions
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
		table.WithFocused(false),
		table.WithHeight(12),
	)

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
	t.SetStyles(focused)

	unfocused = table.DefaultStyles()
	unfocused.Header = unfocused.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	unfocused.Selected = unfocused.Selected.
		Foreground(lipgloss.Color("#bbbbbb"))

	m := model{
		table:          t,
		progress:       progress.New(progress.WithDefaultGradient()),
		textInput:      autocomplete.Init(),
		menus:          model{}.NewMenu(),
		DataSet:        entities.InitDataSet(),
		multiTextInput: textinputs.InitialModel(),
		sender:         sender.NewModel(),
	}

	m.DataSet.Images["testImage1"] = &entities.Image{
		"testImage1.jpg",
		".jpg",
		"F:\\lora",
		0,
		nil,
		entities.Caption{},
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
		for i := 0; i < 10; i++ {
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
