package bubbles

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (

	// Tabs.
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(color).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(color).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type Tab struct {
	Title string
	tea.Model
}

type TabPane struct {
	active int
	tabs   []Tab
}

func NewTabPane(tabs ...Tab) TabPane {
	return TabPane{
		active: 0,
		tabs:   tabs,
	}
}

func (t *TabPane) Next() {
	t.active = (t.active + 1 + len(t.tabs)) % len(t.tabs)
}

func (t *TabPane) Prev() {
	t.active = (t.active - 1 + len(t.tabs)) % len(t.tabs)
}

func (t TabPane) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, tab := range t.tabs {
		cmds = append(cmds, tab.Init())
	}
	return tea.Batch(cmds...)
}

func (t TabPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "pgdown":
			t.Next()
			return t, nil
		case "pgup":
			t.Prev()
			return t, nil
		}
		var cmd tea.Cmd
		t.tabs[t.active].Model, cmd = t.tabs[t.active].Update(msg)
		cmds = append(cmds, cmd)
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft && msg.Y <= 2 {
			(&t).activeByX(msg.X)
		}
	case Broadcast:
		for i, tab := range t.tabs {
			updated, cmd := tab.Model.Update(msg.Msg)
			t.tabs[i].Model = updated
			cmds = append(cmds, cmd)
		}
	default:
		var cmd tea.Cmd
		t.tabs[t.active].Model, cmd = t.tabs[t.active].Update(msg)
		cmds = append(cmds, cmd)
	}

	return t, tea.Batch(cmds...)
}

func (t TabPane) View() string {

	doc := strings.Builder{}

	content := t.tabs[t.active].View()

	tabWidths := t.calculateTabWidths()

	var renderedTabs []string
	//Title Bar
	for i, tab := range t.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(t.tabs)-1, i == t.active
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		}
		if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Width(tabWidths[i]).Render(tab.Title))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(
		windowStyle.
			Width(lipgloss.Width(row) - 2).
			Render(content),
	)
	return docStyle.Render(doc.String())
}

func (t TabPane) calculateTabWidths() []int {

	var widths []int

	content := t.tabs[t.active].View()

	titles := []string{}
	for _, tab := range t.tabs {
		titles = append(titles, tab.Title)
	}

	tabBarWidth := (maxes(titles...) + 2) * len(t.tabs)

	width := max(
		lipgloss.Width(content),
		tabBarWidth,
	)

	for i := range t.tabs {
		isLast := i == len(t.tabs)-1
		var tabWidth int
		if width%len(t.tabs) == 0 {
			tabWidth = width / len(t.tabs)
		} else {
			tabWidth = width / (len(t.tabs) + 1)
			if isLast {
				tabWidth += width % len(t.tabs)
			}
		}
		widths = append(widths, tabWidth)
	}
	return widths
}

func (t *TabPane) activeByX(x int) {
	widths := t.calculateTabWidths()
	var cur int
	for i, w := range widths {
		cur += w + 2
		if x <= cur {
			t.active = i
			return
		}
	}
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxes(elems ...string) int {
	max := lipgloss.Width(elems[0])
	for _, elem := range elems[1:] {
		if cur_len := lipgloss.Width(elem); max < cur_len {
			max = cur_len
		}
	}
	return max
}
