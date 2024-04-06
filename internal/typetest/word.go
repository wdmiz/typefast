package typetest

import (
	"unicode/utf8"
)

type LetterStatus int

const (
	LSNone LetterStatus = iota
	LSCorrect
	LSWrong
	LSOverflow
	LSCursor
)

type Word struct {
	Letters []rune
	Input   []rune
	Cursor  int
}

func NewWord(word string) *Word {
	l := utf8.RuneCountInString(word)
	if l <= 0 {
		return nil
	}

	letters := make([]rune, l)
	progress := make([]rune, l)

	for i, r := range []rune(word) {
		letters[i] = r
		progress[i] = 0
	}

	return &Word{letters, progress, 0}
}

func (w *Word) Runes() ([]rune, []LetterStatus) {
	length := w.Len()
	runes := make([]rune, length)
	state := make([]LetterStatus, length)

	for i := range w.Letters {
		l := w.Input[i]
		switch l {
		case w.Letters[i]:
			runes[i] = l
			state[i] = LSCorrect
		case 0, ' ':
			runes[i] = w.Letters[i]
			state[i] = LSNone
		default:
			runes[i] = l
			state[i] = LSWrong
		}
	}

	for i := len(w.Letters); i < length; i++ {
		runes[i] = w.Input[i]
		state[i] = LSOverflow
	}

	return runes, state
}

func (w *Word) Len() int {
	for i := len(w.Input) - 1; i >= len(w.Letters); i-- {
		if w.Input[i] != 0 {
			return i + 1
		}
	}

	return len(w.Letters)
}
