package king

type CompTest struct {
	Compfile string
	Comptest string
}

type T struct {
	Do   T1 `cmd:"" aliases:"d" help:"do it"`
	More T1 `cmd:"MorethenEver" aliases:"more" help:"do it again" description:"whay more do you want."`

	EvenMore T2 `cmd:"" aliases:"e" help:"do it another time" description:"When running this command you need..."`
}

type T1 struct {
	Status *string `enum:"ok,setup,dst,archive,rm" help:"Set the status for this volume to *STATUS*. See **VOLUME STATUS** section." aliases:"stat" short:"s"`

	Volume string `arg:"" placeholder:"server[:vol]|ID|vol" help:"Volume to change." completion:"echo a b c" type:"existingvolume"`
	Arg    string `arg:"" help:"This is an arg."`
}

type T2 struct {
	Status *string `enum:"ok,setup,dst,archive,rm" help:"Set the status for this volume to *STATUS*. See **VOLUME STATUS** section." aliases:"stat" short:"s"`

	DoEvenMore   T3 `cmd:""`
	WhatEvenMore T4 `cmd:""`
}

type (
	T3 struct{}
	T4 struct{}
)
