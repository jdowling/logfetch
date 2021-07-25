package reversescan

import (
	"bytes"
	"io"
	"log"
	"os"
)

const (
	// TODO: maybe more thought into if this is optimal
	// balance potential system impacts
	// with the current implementation this will also mangle
	// very long individual log lines.
	// so consider buffering based on any system line lenght limits if the system logging service
	// enforces one or do something else to handle that - possible to check with PO about requirements
	// around very long log lines.
	maxBufferSize = 1024
)

/**
 * TL;DR - Read hunks of the file from the end of input back to the beginning of input
 */
type ReverseScanner struct {
	reader    io.ReaderAt // file to read backwards
	position  int64       // start of the next line
	pseudoEOF int64
	buffer    []byte
}

func New(f *os.File) *ReverseScanner {
	// figure out how long the file is using stat.
	stat, err := f.Stat()
	if err != nil {
		// TODO: maybe don't panic here or catch at API layer.
		log.Fatal(err)
	}
	fileLen := stat.Size() - 1

	return &ReverseScanner{reader: f, position: fileLen, pseudoEOF: fileLen}
}

func (s *ReverseScanner) Scan() bool {
	/**
	  So this bit is inefficient but ought to be far more efficient than reading the whole
	  file & we limit the RAM impacts of how much we're going to read at a time.
	  I think it would be more efficient to do fewer reads and buffer more into memory at once.
	  So instead of reading every time we shift lines we could find newlines within the buffer
	  then check to see if the buffer contains a new fill newline or not - if it does just read it an adjust EOF
	  offset within buffer or something. When we run out of new full lines if we're not at start of file
	  then use the current beginning of last line read as the new EOF and read a new full buffer of
	  maxBufferSize to start working through.

	  This was easier to impl within the time window b/c the ptr/offset math was more obvious
	  to me.
	*/
	if s.position > -1 {
		// position will be either the EOF or end of previous line
		s.pseudoEOF = s.position                    // update our logical EOF that we're tracking.
		nextPos := max(0, s.position-maxBufferSize) // don't read before start of file.

		// read in new buffer - don't exceed expected bytes to read.
		expected_read_len := s.pseudoEOF - nextPos
		s.buffer = make([]byte, expected_read_len)
		bytesRead, err := s.reader.ReadAt(s.buffer, nextPos)
		if err != nil {
			// TODO: maybe handle this error more gracefully.
			log.Println("Unexpected error while reading buffer from file. err:", err)
			return false
		}
		if expected_read_len != int64(bytesRead) {
			// TODO: maybe handle this error more gracefully.
			log.Println("Expected to read:", expected_read_len, "bytes but instead read:", bytesRead)
			return false
		}

		// find start of next line looking for \n in buffer.
		// TODO: is there a better way to identify log files? Ex: is there a reasonable cross-distro regex? encoding concerns? etc.
		// this probably doesn't report properly anything where the log line itself contains newline characters.
		lineStartInBuffer := bytes.LastIndexByte(s.buffer, '\n')
		s.position = s.pseudoEOF - (int64(len(s.buffer)) - int64(lineStartInBuffer)) // index inside buffer is kind of like an offset from total file index.
		if s.position > 0 && lineStartInBuffer == -1 {
			log.Fatal(err)
		}
		return true
	} else {
		// we're at the start of the line; nothing left.
		return false
	}
}

// Seems that there's int/floating point issues
// with math.Max type functions in go.
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (s *ReverseScanner) Text() string {
	// at this point s.position should always be start of next line.
	// and s.pseudoEOF should always be the end of it.
	// not counting whitespace.
	readlen := s.pseudoEOF - s.position
	buff := make([]byte, readlen)
	readPos := max(0, s.position)
	s.reader.ReadAt(buff, readPos)
	return string(bytes.TrimSpace(buff))
}
