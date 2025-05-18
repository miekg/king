package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestZsh(t *testing.T) {
	parser := kong.Must(&T{})
	z := &Zsh{}
	z.Completion(parser.Model.Node, "mijnexe")
	if err := z.Write(); err != nil {
		t.Fatal(err)
	}
	parser = kong.Must(&T1{})
	z.Completion(parser.Model.Node, "t1")
	if err := z.Write(); err != nil {
		t.Fatal(err)
	}
}
