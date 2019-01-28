// +build js

package formatter

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
