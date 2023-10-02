package tui

import (
	"github.com/charmbracelet/bubbles/table"
)

var focused = table.DefaultStyles()
var unfocused = table.DefaultStyles()

func (m model) View() string {
	var view string

	view += m.menus[0].Menu.View() + "\n"
	// set styles for each menu
	for menuID, currentMenu := range m.menus {
		if currentMenu.Menu.Focused() {
			m.menus[menuID].Menu.SetStyles(focused)
			if menuID > 0 {
				view += "\n" + m.menus[menuID].Menu.View() + "\n"
			}
		} else {
			m.menus[menuID].Menu.SetStyles(unfocused)
		}
	}

	// TODO: Fix height viewport

	if m.showProgress {
		view += m.progress.View() + "\n"
	}

	view = baseStyle.Render(view) + "\n"

	if m.showTextInput {
		textInputView := m.textInput.View()
		view += baseStyle.Render(textInputView) + "\n"
	}

	if m.showMultiInput {
		textInputView := m.multiTextInput.View()
		view += baseStyle.Render(textInputView) + "\n"
	}

	if m.sender.Active {
		view += m.sender.View() + "\n"
	}

	return view
}
