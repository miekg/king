package king

import (
	"html/template"
	"os"
	"os/exec"
	"testing"

	"github.com/alecthomas/kong"
)

func TestZsh(t *testing.T) {
	parser := kong.Must(&T{})
	manf := &kong.Flag{Value: &kong.Value{Name: "man", Help: "how context-sensitive manual page.", Tag: &kong.Tag{}}}
	z := &Zsh{Flags: []*kong.Flag{manf}}
	z.Completion(parser.Model.Node, "myexe")
	z.Write()
	tmpl, err := template.New("comptest.zsh.tmpl").ParseFiles("comptest.zsh.tmpl")
	if err != nil {
		panic(err)
	}
	c := CompTest{Compfile: "_myexe", Comptest: "myexe -"}

	f, err := os.Create("comptest.zsh")
	err = tmpl.Execute(f, c)
	if err != nil {
		t.Fatal(err)
	}
	os.Chmod("comptest.zsh", 0o755)
	cmd := exec.Command("zsh", "./comptest.zsh")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	println(string(out))
}

func TestZshServerCh(t *testing.T) {
	z := &Zsh{}
	parser := kong.Must(&ServerCh{})
	z.Completion(parser.Model.Node, "ChVolumeServer")
	z.Write()
}
