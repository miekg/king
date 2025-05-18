package king

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestBash(t *testing.T) {
	parser := kong.Must(&T{})
	b := &Bash{}
	b.Completion(parser.Model.Node, "mijnexe")
	b.Write()
	println(string(b.Out()))
}

func TestBashServerCh(t *testing.T) {
	type ServerCh struct {
		Force         *bool   `help:"Force a change. Used in combination with some other flags." short:"f"`
		Commit        int     `help:"Commit percentage for this server. This is an integer from 50 to 200. For home servers the default is 120, and 80 for volume servers."`
		BackupType    *string `help:"Set the backup type for this server." enum:"kopia"`
		BackupBackend *string `help:"Set the backup backend for this server." enum:"s3"`
		BackupConnect string  `help:"How should the backup connect to the backend storage."`
		BackupSecret  string  `help:"Backup encryption key."`
		Status        *string `enum:"ok,setup,dfs,rm," help:"Set the server's status to *STATUS*."`
		Comment       *string `help:"Set a comment. Use the empty string to remove a comment."`
		Owner         *string `help:"Set an owner, this is a free-form string, see ListVolumeServers(1) for a current list."`
		Server        string  `arg:"" placeholder:"server|ID" help:"Server to change." completion:"c volume-server list --comp" type:"existingvolumeserver"`
	}
	b := &Bash{}
	parser := kong.Must(&ServerCh{})
	b.Completion(parser.Model.Node, "ChVolumeServer")
	println(string(b.Out()))
}
