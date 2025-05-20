package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestBash(t *testing.T) {
	parser := kong.Must(&T{})
	b := &Bash{}
	b.Completion(parser.Model.Node, "mijnexe")
	b.Write()
}
