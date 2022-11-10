package tui

import (
	"time"

	"github.com/cspengl/teatime/internal/timetools"
	"github.com/cspengl/teatime/internal/tui/bubbles"

	tea "github.com/charmbracelet/bubbletea"
)

type mode uint8

const (
	Interactive mode = iota
	Clock
	Timer
)

func LoadProgram(m mode, args []string) *tea.Program {

	var root tea.Model
	switch m {
	case Interactive:
		root = interactive()
	case Clock:
		root = bubbles.NewClock()
	case Timer:
		root = bubbles.NewTimer(11 * time.Second)
	default:
		root = interactive()
	}

	return tea.NewProgram(
		model{Model: root},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
}

type model struct {
	tea.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}

type interactiveModel struct {
	tabs bubbles.TabPane

	timer *timetools.Timer
}

func interactive() interactiveModel {
	m := interactiveModel{
		tabs: bubbles.NewTabPane(
			bubbles.Tab{
				Title: "Clock",
				Model: bubbles.NewClock(),
			},
			bubbles.Tab{
				Title: "Timer",
				Model: bubbles.NewTimer(5 * time.Second),
			},
			bubbles.Tab{
				Title: "Break Schedule",
				Model: bubbles.NewBreakScheduler(10*time.Second, 5*time.Second),
			},
		),
		timer: timetools.NewTimer(5 * time.Second),
	}
	return m
}

func (m interactiveModel) Init() tea.Cmd {
	return m.tabs.Init()
}

func (m interactiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	var updated tea.Model
	updated, cmd = m.tabs.Update(msg)
	m.tabs = updated.(bubbles.TabPane)
	return m, cmd
}

func (m interactiveModel) View() string {
	return m.tabs.View()
}
