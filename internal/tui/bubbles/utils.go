package bubbles

import tea "github.com/charmbracelet/bubbletea"

type Broadcast struct {
	tea.Msg
}

func Send(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
