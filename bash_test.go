package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestBash(t *testing.T) {
	parser := kong.Must(&T{})
	b := &Bash{}
	b.Run(parser.Model.Node, "mijnexe")
	// mijnexe.bash
}
