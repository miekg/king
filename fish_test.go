package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func compfTest(t *testing.T, completionfile, exe string) {
	// test actual completion ala zsh
}

func TestFish(t *testing.T) {
	parser := kong.Must(&T{})
	f := &Fish{}
	f.Completion(parser.Model.Node, "myexe")
	println(string(f.Out()))
}
