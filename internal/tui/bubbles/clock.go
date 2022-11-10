package bubbles

import (
	"fmt"
	"strings"
	"time"

	"github.com/cspengl/teatime/internal/timetools"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gen2brain/beeep"

	"github.com/common-nighthawk/go-figure"
)

var (
	color = lipgloss.Color("#9BCD9B")

	defaultBorder = lipgloss.Border{
		Bottom: "─",
	}

	clockElementStyle = lipgloss.NewStyle().
				Height(5).
				Width(18).
				Border(defaultBorder).
				BorderForeground(color).
				Foreground(color)

	selectedElementBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomRight: "┘",
		BottomLeft:  "└",
	}

	selectedClockElementStyle = clockElementStyle.Copy().
					Border(selectedElementBorder, true).
					BorderForeground(color)
)

const (
	timeFormat = "15:04:05"
)

type Clock struct {
	*clockDisplay
}

func secondTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return Broadcast{timeMsg(t)}
	})
}

func NewClock() *Clock {
	c := &Clock{
		clockDisplay: newClockDisplay(),
	}
	c.setTime(time.Now())
	return c
}

func (c *Clock) Init() tea.Cmd {
	return secondTick()
}

type timeMsg time.Time

func (c *Clock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if t, ok := msg.(timeMsg); ok {
		c.setTime(time.Time(t))
		return c, secondTick()
	}
	return c, nil
}

func (c *Clock) setTime(t time.Time) {
	fmt.Sscanf(
		t.Format(timeFormat),
		"%d:%d:%d",
		&c.elems[0],
		&c.elems[1],
		&c.elems[2],
	)
}

type Timer struct {
	*clockDisplay
	timer timetools.Timer
}

func NewTimer(d time.Duration) *Timer {
	t := &Timer{
		clockDisplay: newClockDisplay(),
		timer:        *timetools.NewTimer(d),
	}
	t.selected = seconds
	t.Set(d)
	return t
}

func (t Timer) Get() time.Duration {
	d, _ := time.ParseDuration(
		fmt.Sprintf(
			"%dh%dm%ds",
			t.elems[0],
			t.elems[1],
			t.elems[2],
		),
	)
	return d
}

type durationMsg time.Duration

type timerAction uint8

const (
	TimerStart timerAction = iota
	TimerStop
	TimerReset
)

type TimerDone string

func (t *Timer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			if !t.timer.IsRunning() {
				return t, Send(TimerStart)
			}
			return t, Send(TimerStop)
		case "up":
			if !t.timer.IsRunning() {
				switch t.selected {
				case hours:
					t.Set(t.Get() + time.Hour)
				case minutes:
					t.Set(t.Get() + time.Minute)
				case seconds:
					t.Set(t.Get() + time.Second)
				}
			}
		case "down":
			if !t.timer.IsRunning() {
				switch t.selected {
				case hours:
					t.Set(t.Get() - time.Hour)
				case minutes:
					t.Set(t.Get() - time.Minute)
				case seconds:
					t.Set(t.Get() - time.Second)
				}
			}
		case "tab":
			t.selected = (t.selected + 1) % 3
		case "shift+tab":
			t.selected = (t.selected - 1) % 3
		case "r":
			t.Set(0)
		}
	case TimerDone:
		beeep.Notify("Teatime", string(msg), "")
		return t, Send(TimerReset)
	case timerAction:
		switch msg {
		case TimerStart:
			t.timer.Set(t.Get())
			t.timer.Start()
			return t, Send(durationMsg(t.timer.Get()))
		case TimerStop:
			t.timer.Stop()
		case TimerReset:
			t.timer.Reset()
			t.Set(t.timer.Get())
		}
	case durationMsg:
		t.Set(time.Duration(msg))
		if msg <= 0 {
			return t, Send(TimerDone("Timer done!"))
		}
		return t, func() tea.Msg {
			if _, ok := <-t.timer.C; !ok {
				return Send(TimerReset)
			}
			return durationMsg(t.timer.Get())
		}
	}
	return t, nil
}

func (t *Timer) Set(d time.Duration) {
	t.elems[0] = int64(d.Truncate(time.Hour).Hours()) % 100
	t.elems[1] = int64(d.Truncate(time.Minute).Minutes()) % 60
	t.elems[2] = int64(d.Truncate(time.Second).Seconds()) % 60
}

type selectedClockElement uint8

const (
	hours selectedClockElement = iota
	minutes
	seconds
	none
)

type clockDisplay struct {
	selected selectedClockElement
	elems    [3]int64
}

func newClockDisplay() *clockDisplay {
	return &clockDisplay{
		selected: none,
		elems:    [3]int64{},
	}
}

func (c clockDisplay) Init() tea.Cmd {
	return nil
}

func (c clockDisplay) Update(msg tea.Msg) (clockDisplay, tea.Cmd) {
	return c, nil
}

func (c clockDisplay) View() string {
	var renderedElems []string
	for i, elem := range c.elems {
		if i == int(c.selected) {
			renderedElems = append(renderedElems,
				selectedClockElementStyle.Render(
					centerBlock(
						figure.NewFigure(fmt.Sprintf("%02d", elem), "", false).String(),
						18,
					),
				),
			)
		} else {
			renderedElems = append(renderedElems,
				clockElementStyle.Render(
					centerBlock(
						figure.NewFigure(fmt.Sprintf("%02d", elem), "", false).String(),
						18,
					),
				),
			)
		}
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		renderedElems...,
	)
}

func centerBlock(input string, width int) string {
	res := strings.Builder{}
	lines := strings.Split(input, "\n")

	max := len(lines[0])
	for _, line := range lines[1:] {
		if max < len(line) {
			max = len(line)
		}
	}

	leftPad := (width - max) / 2

	for _, line := range lines {
		if leftPad > 0 {
			res.WriteString(strings.Repeat(" ", leftPad))
		}
		res.WriteString(line + "\n")
	}

	return strings.TrimSuffix(res.String(), "\n")
}
