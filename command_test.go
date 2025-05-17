package king

import (
	"fmt"
	"testing"

	"github.com/alecthomas/kong"
)

type T struct {
	Do   T1 `cmd:"" aliases:"d" help:"do it"`
	More T1 `cmd:"" aliases:"m" help:"do it again"`

	EvenMore T2 `cmd:"" aliases:"e" help:"do it another time"`
}

type T1 struct {
	Status *string `enum:"ok,setup,dst,archive,rm" help:"Set the status for this volume to *STATUS*. See **VOLUME STATUS** section." aliases:"stat" short:"s"`

	Volume string `arg:"" placeholder:"server[:vol]|ID|vol" help:"Volume to change." completion:"echo a b c" type:"existingvolume"`
	Arg    string `arg:"" help:"This is an arg."`
}

type T2 struct {
	DoEvenMore   T3 `cmd:""`
	WhatEvenMore T4 `cmd:""`
}

type (
	T3 struct{}
	T4 struct{}
)

func TestCommands(t *testing.T) {
	parser := kong.Must(&T{})
	cmds := commands(parser.Model.Node)
	for i := range cmds {
		println(cmds[i])
	}
}

func TestNodeForCommand(t *testing.T) {
	parser := kong.Must(&T{})
	node := nodeForCommand(parser.Model.Node, "do")
	fmt.Printf("%v\n", node)
	node = nodeForCommand(parser.Model.Node, "bla")
	fmt.Printf("%v\n", node)
}
