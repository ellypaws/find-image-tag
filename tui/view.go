package tui

import (
	"github.com/charmbracelet/bubbles/table"
)

var focused = table.DefaultStyles()
var unfocused = table.DefaultStyles()

func (m model) View() string {
	var view string
	for i, _ := range m.menu {
		// set the styles depending if they're focused
		if m.menu[i].Focused() {
			m.menu[i].SetStyles(focused)
		} else {
			m.menu[i].SetStyles(unfocused)
		}
		//view += m.menu[i].View() + "\n"
	}
	view += m.menu[0].View() + "\n"
	view += m.menu[1].View() + "\n"

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
	return view
}
