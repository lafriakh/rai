package internal

import (
	"bufio"
	"io"
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

		if err := f(input); err != nil {
			break
		}
	}
}
