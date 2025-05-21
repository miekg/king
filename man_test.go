package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestMan(t *testing.T) {
	parser := kong.Must(&T2{})
	manf := &kong.Flag{Value: &kong.Value{Name: "man", Help: "how context-sensitive manual page.", Tag: &kong.Tag{}}}
	m := &Man{
		Flags:     []*kong.Flag{manf},
		Section:   1,
		Area:      "User Commands",
		WorkGroup: "The hard working team",
	}
	tt := &T{}
	parent := kong.Must(tt.EvenMore)
	m.Manual(parser.Model.Node, parent.Model.Node)
	println(string(m.Out()))
}
