package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func compbTest(t *testing.T, completionfile, exe string) {
	// This can be tested in prolly way:
	//COMP_REPLY=() COMPREPLY="" COMP_WORDS=() COMP_CWORD=1  _myexe_completions ; echo "${COMPREPLY[*]}"
}

func TestBash(t *testing.T) {
	parser := kong.Must(&T{})
	b := &Bash{}
	b.Completion(parser.Model.Node, "myexe")
	b.Write()
}
