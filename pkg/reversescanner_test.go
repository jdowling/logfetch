package reversescan

import (
	"log"
	"os"
	"testing"
)

func TestScanText_findsALine(t *testing.T) {
	file, err := os.Open("test.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := New(file)
	if scanner.Scan() {
		line := scanner.Text()
		expected := "3"
		if line != expected {
			t.Errorf("ReverseScanner returned wrong last line: got %v want %v",
				line, expected)
		}
	}
}

func TestScanText_findsAllLines(t *testing.T) {
	file, err := os.Open("test.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := New(file)
	expected := []string{"3", "2", "1"}
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line != expected[count] {
			t.Errorf("ReverseScanner returned wrong last line: got %v want %v",
				line, expected[count])
		}
		count++
	}
	if count != 3 {
		t.Errorf("ReverseScanner read wrong number of lines: got %v want 3",
			count)
	}
}
