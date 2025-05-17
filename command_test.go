package king

import (
	"fmt"
	"testing"

	"github.com/alecthomas/kong"
)

type T struct {
	Do   T1 `cmd:"" aliases:"d" help:"do it"`
	More T1 `cmd:"" aliases:"m" help:"do it again"`

	EvenMore T2 `cmd:""`
}

type T1 struct {
	Status *string `enum:"ok,setup,dst,archive,rm" help:"Set the status for this volume to *STATUS*. See **VOLUME STATUS** section." aliases:"stat" short:"s"`

	Volume string `arg:"" placeholder:"server[:vol]|ID|vol" help:"Volume to change." completion:"echo a b c" type:"existingvolume"`
	Arg    string `arg:"" help:"This is an arg."`
}

type T2 struct {
	DoEvenMore bool `cmd:""`
}

func TestCommands(t *testing.T) {
	parser := kong.Must(&T{})
	cmds := commands(parser.Model.Node)
	fmt.Printf("%v\n", cmds)
}
