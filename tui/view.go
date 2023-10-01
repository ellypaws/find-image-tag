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
		//view += m.menus[menuID].Menu.View() + "\n"
	}

	// TODO: Fix height viewport
	//view += m.menu[2].View() + "\n"
	//for _, menu := range m.menu {
	//	table := menu.View() + "\n"
	//	view += baseStyle.Render(table) + "\n"
	//}
	//view += m.table.View() + "\n"
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

	//if m.sender.Active {
	//	view += m.sender.View() + "\n"
	//}
	view += m.sender.View() + "\n"

	return view
}
