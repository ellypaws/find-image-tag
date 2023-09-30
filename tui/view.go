package tui

func (m model) View() string {
	view := m.table.View() + "\n"
	if m.showProgress {
		view += m.progress.View() + "\n"
	}
	if m.showTextInput {
		view += m.textInput.View() + "\n"
	}
	return baseStyle.Render(view) + "\n"
}
