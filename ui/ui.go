package ui

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/zulily/stevedore/Godeps/_workspace/src/github.com/mgutz/ansi"
	"github.com/zulily/stevedore/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
)

var (
	blank []byte
)

const (
	defaultWidth = 79
)

func init() {
	width, _, err := terminal.GetSize(0)
	if err != nil || !terminal.IsTerminal(0) {
		width = defaultWidth
	}
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

var (
	taskColor       = ansi.ColorCode("blue+h")
	brightTaskColor = ansi.ColorCode("white+h")
	errColor        = ansi.ColorCode("red")
	brightErrColor  = ansi.ColorCode("red+h")
	warnColor       = ansi.ColorCode("yellow")
	brightWarnColor = ansi.ColorCode("yellow+h")
	infoColor       = ansi.ColorCode("white")
	brightInfoColor = ansi.ColorCode("cyan+h")
	reset           = ansi.ColorCode("reset")
)

// Task formats and prints a message to console in the task color.
func Task(msg string, args ...string) {
	if len(args) == 0 {
		fmt.Println(taskColor + msg + reset)
		return
	}
	colored := colorArgs(args, brightTaskColor, taskColor)
	fmt.Printf(msg+"\n", colored...)
}

// Err formats and prints a message to console in the error color.
func Err(msg string, args ...string) {
	fmt.Println(errColor + msg + reset)
	if len(args) == 0 {
		fmt.Println(errColor, msg, reset)
		return
	}
	colored := colorArgs(args, brightErrColor, errColor)
	fmt.Printf(msg+"\n", colored...)
}

// Warn formats and prints a message to console in the warning color.
func Warn(msg string, args ...string) {
	fmt.Println(warnColor + msg + reset)
	if len(args) == 0 {
		fmt.Println(warnColor, msg, reset)
		return
	}
	colored := colorArgs(args, brightWarnColor, warnColor)
	fmt.Printf(msg+"\n", colored...)
}

// Info formats and prints a message to console in the info color.
func Info(msg string, args ...string) {
	if len(args) == 0 {
		fmt.Println(infoColor + msg + reset)
	} else {
		colored := colorArgs(args, brightInfoColor, infoColor)
		fmt.Printf(msg+"\n", colored...)
	}
}

func colorArgs(args []string, color, reset string) []interface{} {
	var colored []interface{}
	for _, v := range args {
		colored = append(colored, color+v+reset)
	}
	return colored
}
