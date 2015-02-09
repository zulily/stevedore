package ui

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	width int
	blank []byte
)

func init() {
	width, _, _ = terminal.GetSize(0)
	blank = []byte("\r" + strings.Repeat(" ", width) + "\r")
}

// Writer wraps another io.Writer and re-impliments the io.Writer interface,
// stripping any newline characters out of the output and replacing them with
// carriage-returns.
type Writer struct {
	w io.Writer
}

// Wrap creates a new Writer which wraps the given io.Writer.
func Wrap(w io.Writer) Writer {
	return Writer{w: w}
}

func (w Writer) Write(p []byte) (int, error) {
	if _, err := w.w.Write(blank); err != nil {
		return 0, err
	}
	p = bytes.Replace(p, []byte("\r"), []byte("\n"), -1)
	p = bytes.Replace(p, []byte("\n"), []byte("\r"), -1)
	p = bytes.Replace(p, []byte("\t"), []byte(" "), -1)
	return w.w.Write(p)
}
