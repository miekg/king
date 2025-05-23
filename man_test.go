package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

var manf = &kong.Flag{Value: &kong.Value{Name: "man", Help: "how context-sensitive manual page.", Tag: &kong.Tag{}}}

func TestMan(t *testing.T) {
	parser := kong.Must(&T{})
	m := &Man{Flags: []*kong.Flag{manf}, Section: 1, Area: "User Commands", WorkGroup: "The hard working team"}
	m.Manual(parser.Model.Node, "even-more do-even-more", "ListEvenMore", "c")
	println(string(m.Out()))
	m.Write()
}

func TestMan2(t *testing.T) {
	parser := kong.Must(&T{})
	m := &Man{Flags: []*kong.Flag{manf}, Section: 1, Area: "User Commands", WorkGroup: "The hard working team"}
	m.Manual(parser.Model.Node, "even-more", "ListMore", "c")
	println(string(m.Out()))
	m.Write()
}
