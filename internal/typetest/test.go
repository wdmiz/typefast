package typetest

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wdmiz/typefast/internal/stats"
)

var styles = map[letterStatus]lipgloss.Style{
	lsNone: lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")),
	lsCorrect: lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")),
	lsWrong: lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")),
	lsOverflow: lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")),
	lsCursor: lipgloss.NewStyle().Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("15")),
}

type Test struct {
	words       []word
	currentWord int
	Stats       stats.Stats
}

func New(str []string) Test {
	words := make([]word, len(str))

	for i := range words {
		if len(str) <= 0 {
			continue
		}
		words[i] = newWord(str[i])
	}

	return Test{words, 0, stats.New()}
}

func (t *Test) String(width int) (string, string) {
	if width <= 0 {
		return "", ""
	}

	var str strings.Builder
	var aux strings.Builder

	cursorPos := (width - 1) / 2

	startOffset := cursorPos - t.words[t.currentWord].Cursor

	word := t.currentWord

	for w := word; startOffset > 0 && w > 0; w-- {
		startOffset -= t.words[w-1].length() + 1
		word--
	}

	letter := 0

	if startOffset < 0 {
		letter = -(startOffset)
		startOffset = 0
	}

	str.WriteString(strings.Repeat(" ", startOffset))
	aux.WriteString(strings.Repeat(" ", startOffset))

	r, s := t.words[word].runes()

	for i := startOffset; i < cursorPos; i++ {
		if letter >= len(r) {
			if word >= len(t.words)-1 {
				break
			}

			str.WriteString(" ")
			aux.WriteString(" ")

			word++
			r, s = t.words[word].runes()
			letter = 0
			continue
		}

		str.WriteString(styles[s[letter]].Render(string(r[letter])))

		switch s[letter] {
		case lsWrong:
			aux.WriteString(styles[s[letter]].Render("~"))
		case lsOverflow:
			aux.WriteString(styles[s[letter]].Render("#"))
		case lsNone:
			aux.WriteString(styles[s[letter]].Render("`"))
		default:
			aux.WriteString(" ")
		}

		letter++
	}

	aux.WriteString(styles[lsCursor].Render("^"))

	for i := cursorPos; i < width; i++ {
		if letter >= len(r) {
			if word >= len(t.words)-1 {
				break
			}

			str.WriteString(" ")

			word++
			r = t.words[word].Letters
			letter = 0
			continue
		}

		str.WriteString(styles[lsNone].Render(string(r[letter])))
		letter++
	}

	return str.String(), aux.String()
}

func (t *Test) Enter(r rune) {
	w := &t.words[t.currentWord]

	if l := len(w.Input); w.Cursor >= l {
		if l >= 4*len(w.Letters) {
			t.Stats.KeystrokeWrong += 1
			return
		}

		new := make([]rune, l+len(w.Letters))
		copy(new, w.Input)
		w.Input = new
	}

	if w.Cursor < len(w.Letters) {
		if r == w.Letters[w.Cursor] {
			t.Stats.KeystrokeCorrect += 1
		} else {
			t.Stats.KeystrokeWrong += 1
		}
	} else {
		t.Stats.KeystrokeWrong += 1
	}

	w.Input[w.Cursor] = r
	w.Cursor++
}

func (t *Test) Space() (end bool) {
	w := &t.words[t.currentWord]

	if w.Cursor >= len(w.Letters) {
		t.Stats.KeystrokeCorrect += 1
	} else {
		t.Stats.KeystrokeWrong += 1
	}

	if t.currentWord >= len(t.words)-1 {
		return true
	}

	t.currentWord += 1

	return false
}

func (t *Test) Delete() {
	w := &t.words[t.currentWord]

	if w.Cursor >= len(w.Input) {
		w.Cursor = w.length()
	}

	if w.Cursor == 0 {
		if t.currentWord > 0 {
			t.currentWord--
		}

		return
	}

	w.Cursor--
	w.Input[w.Cursor] = 0
}
