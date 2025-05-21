package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestMan(t *testing.T) {
	parser := kong.Must(&T{})
	manf := &kong.Flag{Value: &kong.Value{Name: "man", Help: "how context-sensitive manual page.", Tag: &kong.Tag{}}}
	m := &Man{
		Flags:     []*kong.Flag{manf},
		Section:   1,
		Area:      "User Commands",
		WorkGroup: "The hard working team",
	}
	m.Manual(parser.Model.Node, "more")
	println(string(m.Out()))
}
