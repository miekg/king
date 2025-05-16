package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

type T1 struct {
	Status *string `enum:"ok,setup,dst,archive,rm" help:"Set the status for this volume to *STATUS*. See **VOLUME STATUS** section."`
	Volume string  `arg:"" placeholder:"server[:vol]|ID|vol" help:"Volume to change." completion:"c volume list --comp" type:"existingvolume"`
}

func (t *T1) Run(ctx *kong.Context) error { return nil }

func TestZsh(t *testing.T) {
	parser := kong.Must(&T1{})
	Zsh{}.Completion(parser)
}
