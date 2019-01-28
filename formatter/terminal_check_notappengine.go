// +build !appengine,!js,!windows,!aix

package formatter

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
