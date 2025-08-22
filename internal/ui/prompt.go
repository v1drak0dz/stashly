package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var green = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
var yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
var red = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
var cian = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

// ====================
// === LIST ITEM ===
type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }

// ====================
// === SINGLE SELECT ===
type selectModel struct {
	list     list.Model
	selected string
	done     bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i := m.list.SelectedItem()
			if s, ok := i.(listItem); ok {
				m.selected = string(s)
			}
			m.done = true
			return m, tea.Quit
		case "ctrl+c":
			os.Exit(0)
		}
	}
	return m, cmd
}

func (m selectModel) View() string {
	return m.list.View()
}

func AskSelect(msg string, options []string) (string, error) {
	items := []list.Item{}
	for _, o := range options {
		items = append(items, listItem(o))
	}

	l := list.New(items, list.NewDefaultDelegate(), 50, 15)
	l.Title = msg

	model := selectModel{list: l}
	p := tea.NewProgram(model)
	_, err := p.Run()
	if err != nil {
		return "", err
	}

	return model.selected, nil
}

// ====================
// === MULTI SELECT ===
type multiSelectModel struct {
	options  []string
	cursor   int
	selected map[int]bool
	status   map[string]string // "new", "modified", "deleted"
	done     bool
}

func (m *multiSelectModel) Init() tea.Cmd { return nil }

func (m *multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			os.Exit(0)
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *multiSelectModel) View() string {
	s := "Use ↑ ↓ to move, space to select, enter to confirm:\n"
	for i, choice := range m.options {
		check := " "
		if m.selected[i] || m.cursor == i {
			check = "x"
		}

		// pinta conforme o status
		display := choice
		switch m.status[choice] {
		case "new":
			if m.cursor == i {
				display = cian.Bold(true).Render("[" + check + "] " + choice + " (new)")
			} else {
				display = "[" + check + "] " + choice + green.Render(" (new)")
			}
		case "modified":
			if m.cursor == i {
				display = cian.Bold(true).Render("[" + check + "] " + choice + " (modified)")
			} else {
				display = "[" + check + "] " + choice + yellow.Render(" (modified)")
			}
		case "deleted":
			if m.cursor == i {
				display = cian.Bold(true).Render("[" + check + "] " + choice + " (deleted)")
			} else {
				display = "[" + check + "] " + choice + red.Render(" (deleted)")
			}
		}

		// s += fmt.Sprintf("%s [%s] %s\n", cursor, check, display)
		s += fmt.Sprintf("%s\n", display)
	}
	return s
}

func AskMultiSelectColored(msg string, options []string, status map[string]string) ([]string, error) {
	m := &multiSelectModel{
		options:  options,
		cursor:   0,
		selected: make(map[int]bool),
		status:   status,
	}
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		return nil, err
	}

	var selected []string
	for i, ok := range m.selected {
		if ok {
			selected = append(selected, m.options[i])
		}
	}
	return selected, nil
}

// ====================
// === INPUT ===
type inputModel struct {
	textInput textinput.Model
	done      bool
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c":
			os.Exit(0)
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf("%s\n\nPress enter to confirm\n", m.textInput.View())
}

func AskInput(msg string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = msg
	ti.Focus()

	model := inputModel{textInput: ti}
	p := tea.NewProgram(model)
	_, err := p.Run()
	if err != nil {
		return "", err
	}

	return model.textInput.Value(), nil
}

// ====================
// === CONFIRM ===
func AskConfirm(msg string) (bool, error) {
	fmt.Printf("%s [y/N]: ", msg)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}
	return response == "y" || response == "Y", nil
}
