package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner spinner.Model
	done    bool
	result  string
	err     error
}

type resultMsg struct {
	output string
	err    error
}

func runCommand(args []string) tea.Cmd {
	return func() tea.Msg {
		if len(args) == 0 {
			return resultMsg{"", nil}
		}
		cmd := exec.Command(args[0], args[1:]...)
		out, err := cmd.CombinedOutput()
		return resultMsg{string(out), err}
	}
}

func initialModel(args []string) model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: sp}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, runCommand(os.Args[1:]))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.done {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case resultMsg:
		m.done = true
		m.result = msg.output
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
		content := m.result
		if err == nil {
			if rendered, rerr := renderer.Render(m.result); rerr == nil {
				content = rendered
			}
		}
		if m.err != nil {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(content)
		}
		return content
	}
	return fmt.Sprintf("%s Running...", m.spinner.View())
}

func main() {
	m := initialModel(os.Args[1:])
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
