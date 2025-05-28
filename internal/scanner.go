package internal

import (
	"bufio"
	"io"
	"strings"
)

type Scanner struct {
	scanner *bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{scanner: bufio.NewScanner(r)}
}

func (s *Scanner) Scan(f func(input string) error) {
	for {
		if !s.scanner.Scan() {
			break
		}

		input := s.scanner.Text()
		if input == "" {
			continue
		}
		if strings.HasSuffix(input, "/vim") {
			withoutSuffix, _ := strings.CutSuffix(input, "/vim")
			vimOutput, err := EditStringInVim(withoutSuffix)
			if err != nil {
				continue
			}
			input = vimOutput
		}

		if err := f(input); err != nil {
			break
		}
	}
}
