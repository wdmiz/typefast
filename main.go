package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
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
	stopwatch   stopwatch.Model
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

	p := tea.NewProgram(initialModel(test))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}

func initialModel(t *typetest.Test) model {
	return model{t, "", "", 0, stopwatch.NewWithInterval(time.Millisecond), false, false}
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

			if !m.stopwatch.Running() {
				m.running = true
				m.displayText, m.auxText = m.test.String(m.width)
				return m, m.stopwatch.Start()
			}
		}

		m.displayText, m.auxText = m.test.String(m.width)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.displayText, m.auxText = m.test.String(m.width)
	}

	var cmd tea.Cmd

	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var wpmr int
	if s := m.stopwatch.Elapsed().Minutes(); s > 0 {
		wpmr = int(float64(m.test.Stats.KeystrokeCorrect/5) / s)
	}

	var acc float64
	if c, w := m.test.Stats.KeystrokeCorrect, m.test.Stats.KeystrokeWrong; c > 0 {
		acc = float64(c) / float64(c+w)
	}

	wpm := int(float64(wpmr) * acc)

	elapsed := m.stopwatch.Elapsed().Milliseconds()

	timer := fmt.Sprintf("%02d:%02d.%03d", elapsed/60000, elapsed/1000, elapsed%1000)

	if m.quitting {
		return fmt.Sprintf("T:%s, WPM(r):%d(%d), ACC:%d%%, KS(all/ok/wrong):%d/%d/%d\n",
			timer,
			wpm,
			wpmr,
			int(acc*100),
			m.test.Stats.KeystrokeCorrect+m.test.Stats.KeystrokeWrong,
			m.test.Stats.KeystrokeCorrect,
			m.test.Stats.KeystrokeWrong)
	}

	header := timer + " " + fmt.Sprintf("WPM: %d", wpm) + " " + fmt.Sprintf("ACC: %d%%", int(acc*100))

	return headerStyle.Width(m.width).Render(header) +
		"\n" +
		m.displayText +
		"\n" +
		m.auxText +
		"\n" +
		footerStyle.Width(m.width).Render("Press [ESC] to exit")
}
