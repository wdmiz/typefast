package text

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadText(file string) ([]string, error) {
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
		return nil, fmt.Errorf("text is empty")
	}

	return words, nil
}
