package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/wdmiz/gotypefast/internal/text"
	"github.com/wdmiz/gotypefast/internal/typetest"
)

type model struct {
	test        *typetest.Test
	displayText string
	auxText     string
	width       int
	timeStart   time.Time
	running     bool
	quitting    bool
}

var headerStyle lipgloss.Style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("11")).
	Bold(true)

var footerStyle lipgloss.Style = lipgloss.NewStyle().
	Align(lipgloss.Left).
	Foreground(lipgloss.Color("8"))

func main() {
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
			words = make([]string, *wordCountFlag)
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

	test := typetest.NewTest(&words)

	p := tea.NewProgram(initialModel(test), tea.WithFPS(60))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}

func initialModel(t *typetest.Test) model {
	return model{t, "", "", 0, time.Now(), false, false}
}

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
				m.timeStart = time.Now()
			}
		}

		m.displayText, m.auxText = m.test.String(m.width)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.displayText, m.auxText = m.test.String(m.width)
	}

	var cmd tea.Cmd

	tick := func() tea.Msg {
		if m.running {
			return cmd
		}

		return nil
	}

	return m, tick
}

func (m model) View() string {
	var t time.Duration
	if m.running {
		t = time.Since(m.timeStart)
	}

	timer := fmt.Sprintf("%01d:%02d.%01d", int(t.Minutes()), int(t.Seconds())%60, (t.Milliseconds()%1000)/100)

	if m.quitting {
		return fmt.Sprintf("T:%s, WPM(r):%d(%d), ACC:%d%%, KS(all/ok/wrong):%d/%d/%d\n",
			timer,
			int(m.test.Stats.WPM(t)),
			int(m.test.Stats.WPMRaw(t)),
			int(m.test.Stats.Accuracy()*100),
			m.test.Stats.KeystrokeCorrect+m.test.Stats.KeystrokeWrong,
			m.test.Stats.KeystrokeCorrect,
			m.test.Stats.KeystrokeWrong)
	}

	header := timer + " " + fmt.Sprintf("WPM: %d", int(m.test.Stats.WPM(t))) + " " + fmt.Sprintf("ACC: %d%%", int(m.test.Stats.Accuracy()*100))

	return headerStyle.Width(m.width).Render(header) +
		"\n" +
		m.displayText +
		"\n" +
		m.auxText +
		"\n" +
		footerStyle.Width(m.width).Render("Press [ESC] to exit")
}
