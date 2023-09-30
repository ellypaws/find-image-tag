package tui

func (m model) View() string {
	var view string
	view = m.menu[0].View() + "\n"
	//for _, menu := range m.menu {
	//	table := menu.View() + "\n"
	//	view += baseStyle.Render(table) + "\n"
	//}
	//view += m.table.View() + "\n"
	//if m.showProgress {
	//	view += m.progress.View() + "\n"
	//}

	view = baseStyle.Render(view) + "\n"

	if m.showTextInput {
		textInputView := m.textInput.View()
		view += baseStyle.Render(textInputView) + "\n"
	}
	return view
}
