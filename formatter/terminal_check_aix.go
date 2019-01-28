// +build !appengine,!js,!windows,aix

package formatter

import "io"

func checkIfTerminal(w io.Writer) bool {
	return false
}
