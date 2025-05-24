package king

import (
	"fmt"
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

func TestManNoSuchNode(t *testing.T) {
	parser := kong.Must(&T{})
	m := &Man{Section: 1, Area: "User Commands", WorkGroup: "The hard working team"}
	fmt.Printf("%+v\n", parser.Model.Node)
	fmt.Printf("%+v\n", parser.Model.Tag.Get("description"))
	m.Manual(parser.Model.Node, "does-not-exist", "", "c")
}

type WrapT struct {
	Wrap T `cmd:"" help:"my help" description:"my desc"`
}

func TestManMain(t *testing.T) {
	parser := kong.Must(&WrapT{})
	m := &Man{Section: 1, Area: "User Commands", WorkGroup: "The hard working team"}
	m.Manual(parser.Model.Node, "_wrap", "MyExec", "")
	println(string(m.Out()))
}
