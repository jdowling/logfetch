package reversescan

import (
	"bufio"
	"os"
)

type ReverseScanner struct {
	lines []string // TODO: obviously wrong for large files
	index int
}

func New(f *os.File) *ReverseScanner {
	lines := make([]string, 0, 3)
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	// reverse it
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return &ReverseScanner{lines, 0}
}

func (s *ReverseScanner) Scan() bool {
	return s.index != len(s.lines)
}

func (s *ReverseScanner) Text() string {
	line := s.lines[s.index]
	s.index++
	return line
}
