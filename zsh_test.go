package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestZsh(t *testing.T) {
	parser := kong.Must(&T{})
	manf := &kong.Flag{Value: &kong.Value{Name: "man", Help: "how context-sensitive manual page."}}
	z := &Zsh{Flags: []*kong.Flag{manf}}
	z.Completion(parser.Model.Node, "mijnexe")
	println(string(z.Out()))
	z.Write()
}
