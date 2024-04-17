package typetest

import (
	"unicode/utf8"
)

type letterStatus int

const (
	lsNone letterStatus = iota
	lsCorrect
	lsWrong
	lsOverflow
	lsCursor
)

type word struct {
	Letters []rune
	Input   []rune
	Cursor  int
}

func newWord(s string) word {
	l := utf8.RuneCountInString(s)

	if l == 0 {
		return word{nil, nil, 0}
	}

	letters := make([]rune, l)
	input := make([]rune, l+3)

	for i, r := range []rune(s) {
		letters[i] = r
		input[i] = 0
	}

	return word{letters, input, 0}
}

func (w *word) runes() ([]rune, []letterStatus) {
	length := w.length()
	runes := make([]rune, length)
	state := make([]letterStatus, length)

	for i := range w.Letters {
		l := w.Input[i]
		switch l {
		case w.Letters[i]:
			runes[i] = l
			state[i] = lsCorrect
		case 0, ' ':
			runes[i] = w.Letters[i]
			state[i] = lsNone
		default:
			runes[i] = l
			state[i] = lsWrong
		}
	}

	for i := len(w.Letters); i < length; i++ {
		runes[i] = w.Input[i]
		state[i] = lsOverflow
	}

	return runes, state
}

func (w *word) length() int {
	for i := len(w.Input) - 1; i >= len(w.Letters); i-- {
		if w.Input[i] != 0 {
			return i + 1
		}
	}

	return len(w.Letters)
}
