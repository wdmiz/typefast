package typetest

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var styles = map[LetterStatus]lipgloss.Style{
	LSNone: lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")),
	LSCorrect: lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")),
	LSWrong: lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")),
	LSOverflow: lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")),
	LSCursor: lipgloss.NewStyle().Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("15")),
}

type Test struct {
	Words       []*Word
	CurrentWord int
	Stats       *Stats
}

func NewTest(str *[]string) *Test {
	words := make([]*Word, len(*str))

	for i := range words {
		words[i] = NewWord((*str)[i])
	}

	return &Test{words, 0, NewStats()}
}

func (t *Test) String(width int) (string, string) {
	if width <= 0 {
		return "", ""
	}

	var str strings.Builder
	var aux strings.Builder

	cursorPos := (width - 1) / 2

	startOffset := cursorPos - t.Words[t.CurrentWord].Cursor

	word := t.CurrentWord

	for w := word; startOffset > 0 && w > 0; w-- {
		startOffset -= t.Words[w-1].Len() + 1
		word--
	}

	letter := 0

	if startOffset < 0 {
		letter = -(startOffset)
		startOffset = 0
	}

	str.WriteString(strings.Repeat(" ", startOffset))
	aux.WriteString(strings.Repeat(" ", startOffset))

	r, s := t.Words[word].Runes()

	for i := startOffset; i < cursorPos; i++ {
		if letter >= len(r) {
			if word >= len(t.Words)-1 {
				break
			}

			str.WriteString(" ")
			aux.WriteString(" ")

			word++
			r, s = t.Words[word].Runes()
			letter = 0
			continue
		}

		str.WriteString(styles[s[letter]].Render(string(r[letter])))

		switch s[letter] {
		case LSWrong:
			aux.WriteString(styles[s[letter]].Render("~"))
		case LSOverflow:
			aux.WriteString(styles[s[letter]].Render("#"))
		case LSNone:
			aux.WriteString(styles[s[letter]].Render("`"))
		default:
			aux.WriteString(" ")
		}

		letter++
	}

	aux.WriteString(styles[LSCursor].Render("^"))

	for i := cursorPos; i < width; i++ {
		if letter >= len(r) {
			if word >= len(t.Words)-1 {
				break
			}

			str.WriteString(" ")

			word++
			r = t.Words[word].Letters
			letter = 0
			continue
		}

		str.WriteString(styles[LSNone].Render(string(r[letter])))
		letter++
	}

	return str.String(), aux.String()
}

func (t *Test) Enter(r rune) bool {
	w := t.Words[t.CurrentWord]

	if l := len(w.Input); w.Cursor >= l {
		if l >= 4*len(w.Letters) {
			return false
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
	}

	w.Input[w.Cursor] = r
	w.Cursor++

	return false
}

func (t *Test) Space() (end bool) {
	w := t.Words[t.CurrentWord]

	if w.Cursor >= len(w.Input) {
		t.Stats.KeystrokeCorrect += 1
	} else {
		t.Stats.KeystrokeWrong += 1
	}

	if t.CurrentWord >= len(t.Words)-1 {
		return true
	}

	t.CurrentWord += 1

	return false
}

func (t *Test) Delete() {
	w := t.Words[t.CurrentWord]

	if w.Cursor >= len(w.Input) {
		w.Cursor = w.Len()
	}

	if w.Cursor == 0 {
		if t.CurrentWord > 0 {
			t.CurrentWord--
		}

		return
	}

	w.Cursor--
	w.Input[w.Cursor] = 0
}
