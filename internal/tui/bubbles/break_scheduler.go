package bubbles

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
)

type editField uint8

var (
	editStyle = lipgloss.NewStyle().Underline(true)
)

const (
	editWorkTime editField = iota
	editBreakTime
)

type BreakScheduler struct {
	*Timer
	selectedField       editField
	workMode, editMode  bool
	workTime, breakTime time.Duration
}

func NewBreakScheduler(workTime, breakTime time.Duration) *BreakScheduler {
	return &BreakScheduler{
		Timer:         NewTimer(workTime),
		selectedField: editWorkTime,
		editMode:      false,
		workMode:      true,
		workTime:      workTime,
		breakTime:     breakTime,
	}
}

func (bs BreakScheduler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TimerDone:
		if bs.workMode {
			bs.workMode = false
			bs.Set(bs.breakTime)
			beeep.Notify("Teatime", "Time for a break!", "")
		} else {
			bs.workMode = true
			bs.Set(bs.workTime)
			beeep.Notify("Teatime", "Break's over - Lets get back to work", "")
		}
		return bs, Send(TimerStart)
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			bs.editMode = !bs.editMode
			if bs.editMode {
				if bs.selectedField == editWorkTime {
					bs.Set(bs.workTime)
				} else {
					bs.Set(bs.breakTime)
				}
			} else {
				bs.Set(bs.workTime)
			}
		case "r":
			if bs.workMode {
				bs.timer.Set(bs.workTime)
			} else {
				bs.timer.Set(bs.breakTime)
			}
			return bs, nil
		case "enter":
			if bs.editMode {
				if bs.selectedField == editWorkTime {
					bs.workTime = bs.Get()
				} else {
					bs.breakTime = bs.Get()
				}
			}
		case "left":
			if bs.editMode {
				bs.selectedField = editWorkTime
				bs.Set(bs.workTime)
			}
		case "right":
			if bs.editMode {
				bs.selectedField = editBreakTime
				bs.Set(bs.breakTime)
			}
		case "up", "down":
			if !bs.editMode {
				return bs, nil
			}
		}
	}
	updated, cmd := bs.Timer.Update(msg)
	bs.Timer = updated.(*Timer)
	return bs, cmd
}

func (bs BreakScheduler) View() string {

	doc := strings.Builder{}

	var (
		workDuration  = fmt.Sprintf("Work Duration: %s", bs.workTime.String())
		breakDuration = fmt.Sprintf("Break Duration: %s", bs.breakTime.String())
	)
	if bs.editMode {
		if bs.selectedField == editBreakTime {
			breakDuration = editStyle.Render(breakDuration)
		} else {
			workDuration = editStyle.Render(workDuration)
		}
	}
	settings := lipgloss.JoinHorizontal(
		lipgloss.Left, workDuration, "\t", breakDuration,
	)

	doc.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		settings,
		bs.Timer.View(),
	))

	return doc.String()
}
