package king

type CompTest struct {
	Compfile string
	Comptest string
}

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

type ServerCh struct {
	Force         *bool   `help:"Force a change. Used in combination with some other flags." short:"f"`
	Commit        int     `help:"Commit percentage for this server. This is an integer from 50 to 200. For home servers the default is 120, and 80 for volume servers."`
	BackupType    *string `help:"Set the backup type for this server." enum:"kopia"`
	BackupBackend *string `help:"Set the backup backend for this server." enum:"s3"`
	BackupConnect string  `help:"How should the backup connect to the backend storage."`
	BackupSecret  string  `help:"Backup encryption key." env:"USERDB_TOTP_KEY"`
	Status        *string `enum:"ok,setup,dfs,rm," help:"Set the server's status to *STATUS*."`
	Comment       *string `help:"Set a comment. Use the empty string to remove a comment."`
	Owner         *string `help:"Set an owner, this is a free-form string, see ListVolumeServers(1) for a current list."`
	Server        string  `arg:"" placeholder:"server|ID" help:"Server to change." completion:"echo volume-server list --comp" type:"existingvolumeserver"`
	Server2       string  `arg:"" placeholder:"server|ID" help:"Server to change." completion:"echo volume-server list --comp" type:"existingvolumeserver"`
}
