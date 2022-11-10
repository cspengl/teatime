package bubbles

import tea "github.com/charmbracelet/bubbletea"

type TextPane struct {
	Text string
}

func (t TextPane) Init() tea.Cmd {
	return nil
}

func (t TextPane) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t TextPane) View() string {
	return t.Text
}
