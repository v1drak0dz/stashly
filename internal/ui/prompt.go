package ui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
// === MULTI SELECT COM DIFF + SCROLL ===
type multiSelectModel struct {
	options  []string
	cursor   int
	offset   int // scroll da lista da esquerda
	selected map[int]bool
	status   map[string]string // "new", "modified", "deleted"

	viewport viewport.Model // painel direito (diff)
	diff     string

	height int
	width  int

	done bool
}

func (m *multiSelectModel) Init() tea.Cmd {
	// define tamanho do viewport do diff
	m.height = 20
	m.width = 80
	m.viewport = viewport.New(m.width, m.height)
	return m.loadDiff()
}

func (m *multiSelectModel) loadDiff() tea.Cmd {
	if len(m.options) == 0 {
		m.diff = ""
		m.viewport.SetContent("")
		return nil
	}
	file := m.options[m.cursor]

	// executa git diff --color
	out, _ := exec.Command("git", "diff", "--color=always", file).CombinedOutput()
	m.diff = string(out)
	if m.diff == "" {
		m.diff = "(sem alterações para mostrar)"
	}
	m.viewport.SetContent(m.diff)
	return nil
}

func (m *multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// pega largura e altura do terminal
		m.width = msg.Width - 5   // reserva espaço para borda
		m.height = msg.Height - 2 // reserva header/margem
		m.viewport.Width = m.width
		m.viewport.Height = m.height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			os.Exit(0)

		// navegação na lista de arquivos
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset-- // scroll pra cima
				}
				return m, m.loadDiff()
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
				if m.cursor >= m.offset+m.height-2 { // "-2" p/ header
					m.offset++ // scroll pra baixo
				}
				return m, m.loadDiff()
			}

		// seleção
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]

		// confirma
		case "enter":
			m.done = true
			return m, tea.Quit

		// scroll no diff
		case "pgdown", "ctrl+d":
			m.viewport.ScrollDown(5)
		case "pgup", "ctrl+u":
			m.viewport.ScrollUp(5)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *multiSelectModel) View() string {
	// painel esquerdo: lista de arquivos com scroll
	left := "Use ↑ ↓ para navegar, espaço para marcar, enter para confirmar:\n\n"

	end := m.offset + m.height - 2 // limite visível (reserva header + margem)
	if end > len(m.options) {
		end = len(m.options)
	}

	for i := m.offset; i < end; i++ {
		choice := m.options[i]

		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		check := " "
		if m.selected[i] {
			check = "x"
		}

		display := fmt.Sprintf("[%s] %s", check, choice)
		switch m.status[choice] {
		case "new":
			display = green.Render(display + " (new)")
		case "modified":
			display = yellow.Render(display + " (modified)")
		case "deleted":
			display = red.Render(display + " (deleted)")
		}

		if i == m.cursor {
			display = cian.Bold(true).Render(cursor + " " + display)
		} else {
			display = cursor + " " + display
		}
		left += display + "\n"
	}

	// painel direito: diff com scroll
	right := m.viewport.View()

	// junta em colunas lado a lado
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(45).Height(m.height).Border(lipgloss.NormalBorder()).Render(left),
		lipgloss.NewStyle().Width(m.width - 45).Height(m.height).Border(lipgloss.NormalBorder()).Render(right),
	)
}

func AskMultiSelectColored(msg string, options []string, status map[string]string) ([]string, error) {
	m := &multiSelectModel{
		options:  options,
		cursor:   0,
		offset:   0,
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
