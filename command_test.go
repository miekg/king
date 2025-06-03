package king

import "time"

type CompTest struct {
	Compfile string
	Comptest string
}

type T struct {
	Do   T1 `cmd:"" aliases:"d" help:"do it"`
	More T1 `cmd:"MorethenEver" aliases:"again" help:"do it again" description:"T1 whay more do you want."`

	EvenMore T2 `cmd:"" aliases:"more" help:"do it another time" description:"T2 When running this command you need..."`
}

type T1 struct {
	Status      *string    `placeholder:"status" enum:"ok,setup,dst,archive,rm" help:"Set the status for this volume to *STATUS*. See **VOLUME STATUS** section." aliases:"stat" short:"s"`
	Enddate     *time.Time `help:"Set the end date." format:"2006-01-02" aliases:"afloopdatum" group:"end date"`
	File        string     `help:"complete this file" completion:"<file>"`
	SuperString string     `help:"complete this string" completion:"echo bla bloep"`

	Volume string `arg:"" placeholder:"server[:vol]|ID|vol" help:"Volume to change." completion:"echo a b c" type:"existingvolume"`
	Arg    string `arg:"" help:"This is an arg."`
}

type T2 struct {
	Status *string `placeholder:"status" enum:"ok,setup,dst,archive,rm" help:"Set the status for this volume to *STATUS*. See **VOLUME STATUS** section." aliases:"stat" short:"s"`

	DoEvenMore   T3 `cmd:"" help:"do it agian, but more." description:"T3: this is the thing we want to see."`
	WhatEvenMore T4 `cmd:"" help:"do it again, but even more."`
}

type (
	T3 struct {
		Bool   *bool   `help:"allow a bool." negetable:""`
		String *string `help:"allow a string." completion:"<file>"`

		Arg string `arg:"" help:"This is another arg." placeholder:"bliep"`
	}
	T4 struct{}
)
