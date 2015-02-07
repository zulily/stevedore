package ui

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	width int
)

func init() {
	width, _, _ = terminal.GetSize(0)
}

func foobar() {
	fmt.Printf("%d%s", 123)
}

// This is a terrible comment, and should be caught by golint
// This line won't help either.
type Writer struct {
	w      io.Writer
	length int
}

func Wrap(w io.Writer) Writer {
	return Writer{w: w}
}

func (w Writer) Write(p []byte) (int, error) {
	if _, err := w.w.Write([]byte("\r" + strings.Repeat(" ", width) + "\r")); err != nil {
		return 0, err
	}
	p = bytes.Replace(p, []byte("\r"), []byte("\n"), -1)
	p = bytes.Replace(p, []byte("\n"), []byte("\r"), -1)
	p = bytes.Replace(p, []byte("\t"), []byte(" "), -1)
	return w.w.Write(p)
}
