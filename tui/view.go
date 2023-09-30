package tui

func (m model) View() string {
	view := m.table.View() + "\n"
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
