package text

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

type Randomizer struct {
	words []string
}

func NewRandomizer(file string) (*Randomizer, error) {
	pathAbs, err := filepath.Abs(file)

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(pathAbs)

	if err != nil {
		return nil, err
	}

	words := strings.Fields(string(data))

	if len(words) == 0 {
		return nil, fmt.Errorf("dictionary is empty")
	}

	return &Randomizer{words}, nil
}

func (r *Randomizer) Word() string {
	return r.words[rand.Intn(len(r.words))]
}
