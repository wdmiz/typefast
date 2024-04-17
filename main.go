package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/wdmiz/typefast/internal/text"
	"github.com/wdmiz/typefast/internal/typetest"
)

// Styles
var headerStyle lipgloss.Style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("11")).
	Bold(true)

var footerStyle lipgloss.Style = lipgloss.NewStyle().
	Align(lipgloss.Left).
	Foreground(lipgloss.Color("8"))

type model struct {
	test      typetest.Test
	startTime time.Time

	displayText string
	auxText     string

	width int

	running  bool
	quitting bool
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(
		time.Millisecond*100,
		func(_ time.Time) tea.Msg {
			return tickMsg{}
		},
	)
}

func main() {
	// Run flags
	dictFlag := flag.String("dict", "", "Path to file containg words to generate test text from")
	textFlag := flag.String("text", "", "Path to file containg test text")
	wordCountFlag := flag.Int("words", 100, "Number of words in text")
	flag.Parse()

	var words []string
	var err error

	if *textFlag != "" {
		words, err = text.LoadText(*textFlag)

	} else if *dictFlag != "" {
		var wordRand *text.Randomizer
		wordRand, err = text.NewRandomizer(*dictFlag)

		if err == nil {
			wordCount := *wordCountFlag
			if wordCount > 2<<15 {
				wordCount = 2 << 15
			}

			words = make([]string, wordCount)
			for w := range words {
				words[w] = wordRand.Word()
			}
		}

	} else {
		err = fmt.Errorf("no text or dictionary provided")
	}

	if err != nil {
		fmt.Printf("Failed to prepare test text: %v\n", err.Error())
		os.Exit(1)
	}

	test := typetest.New(words)

	p := tea.NewProgram(initModel(test), tea.WithFPS(60))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}

func initModel(t typetest.Test) model {
	return model{
		test:      t,
		startTime: time.Now(),
	}
}

// Bubbletea model functions

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyBackspace:
			m.test.Delete()

		case tea.KeySpace:
			if m.running && m.test.Space() {
				m.quitting = true
				return m, tea.Quit
			}

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				m.test.Enter(r)
			}

			if !m.running {
				m.running = true
				m.startTime = time.Now()
				m.displayText, m.auxText = m.test.String(m.width)
				return m, tick()
			}
		}

		m.displayText, m.auxText = m.test.String(m.width)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.displayText, m.auxText = m.test.String(m.width)

	case tickMsg:
		if m.running {
			return m, tick()
		}
	}

	return m, nil
}

func (m model) View() string {
	t := time.Since(m.startTime)

	stopwatch := fmt.Sprintf(
		"%01d:%02d.%01d",
		int(t.Minutes()),
		int(t.Seconds())%60,
		(t.Milliseconds()%1000)/100,
	)

	wpm := int(m.test.Stats.WPM(t))
	wpmr := int(m.test.Stats.WPM(t))
	acc := int(m.test.Stats.Accuracy() * 100)

	if m.quitting {
		return fmt.Sprintf("T:%s, WPM(r):%d(%d), ACC:%d%%, KS(all/ok/wrong):%d/%d/%d\n",
			stopwatch,
			wpm,
			wpmr,
			acc,
			m.test.Stats.KeystrokeCorrect+m.test.Stats.KeystrokeWrong,
			m.test.Stats.KeystrokeCorrect,
			m.test.Stats.KeystrokeWrong)
	}

	header := fmt.Sprintf("%s WPM: %d ACC: %d%% ", stopwatch, wpm, acc)

	return headerStyle.Width(m.width).Render(header) +
		"\n" +
		m.displayText +
		"\n" +
		m.auxText +
		"\n" +
		footerStyle.Width(m.width).Render("Press [ESC] to exit")
}
