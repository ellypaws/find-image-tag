package sender

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"os"
	"strings"
	"time"
)

// A simple example that shows how to send messages to a Bubble Tea program
// from outside the program using Program.Send(Msg).

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type ResultMsg struct {
	Duration time.Duration
	Food     string
}

func (r ResultMsg) String() string {
	if r.Duration == 0 {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("🍔 Ate %s %s", r.Food,
		durationStyle.Render(r.Duration.String()))
}

type Model struct {
	Spinner spinner.Model
	Results []ResultMsg
	Active  bool
}

func NewModel() Model {
	const numLastResults = 5
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return Model{
		Spinner: s,
		Results: make([]ResultMsg, numLastResults),
	}
}

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//case tea.KeyMsg:
	//	m.Active = false
	//	return m, tea.Quit
	case ResultMsg:
		m.Results = append(m.Results[1:], msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) View() string {
	var s string

	if m.Active {
		s += m.Spinner.View() + " Eating food..."
	} else {
		s += "That’s all for today!"
	}

	s += "\n\n"

	for _, res := range m.Results {
		s += res.String() + "\n"
	}

	//if m.Active {
	//	s += helpStyle.Render("Press any key to exit")
	//}

	if !m.Active {
		s += "\n"
	}

	return appStyle.Render(s)
}

func main() {
	p := tea.NewProgram(NewModel())

	// Simulate activity
	go func() {
		for {
			pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec
			time.Sleep(pause)

			// Send the Bubble Tea program a message from outside the
			// tea.Program. This will block until it is ready to receive
			// messages.
			p.Send(ResultMsg{Food: RandomFood(), Duration: pause})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func RandomFood() string {
	food := []string{
		"an apple", "a pear", "a gherkin", "a party gherkin",
		"a kohlrabi", "some spaghetti", "tacos", "a currywurst", "some curry",
		"a sandwich", "some peanut butter", "some cashews", "some ramen",
	}
	return food[rand.Intn(len(food))] // nolint:gosec
}
