package king

import (
	"bytes"
	"html/template"
	"os"
	"os/exec"
	"testing"

	"github.com/alecthomas/kong"
)

func compzTest(t *testing.T, completionfile, exe string) []byte {
	tmpl, err := template.New("comptest.zsh.tmpl").ParseFiles("testdata/comptest.zsh.tmpl")
	if err != nil {
		panic(err)
	}
	c := CompTest{Compfile: completionfile, Comptest: exe}

	f, err := os.Create("comptest.zsh")
	err = tmpl.Execute(f, c)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("comptest.zsh")
	os.Chmod("comptest.zsh", 0o755)
	cmd := exec.Command("zsh", "./comptest.zsh")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	return bytes.TrimSpace(out)
}

func TestZsh(t *testing.T) {
	parser := kong.Must(&T{})
	manf := &kong.Flag{Value: &kong.Value{Name: "man", Help: "how context-sensitive manual page.", Tag: &kong.Tag{}}}
	z := &Zsh{Flags: []*kong.Flag{manf}}
	z.Completion(parser.Model.Node, "myexe")
	z.Write()

	tests := []struct {
		exe    string
		expect string
	}{
		{"myexe --", "--help\r\n--man"},
		{"myexe --m", "--man"},
		{"myexe ", "d\r\ndo\r\nagain\r\nmore\r\neven-more"},
	}

	for i := range tests {
		out := compzTest(t, "_myexe", tests[i].exe)
		if string(out) != tests[i].expect {
			t.Errorf("test %d, expected %q, got %q", i, tests[i].expect, string(out))
		}
	}
}

func TestAction(t *testing.T) {
	parser := kong.Must(&T1{})
	z := &Zsh{}
	const exp = `--super-string=[complete this string]:complete this string:_values 'super-string' $(echo bla bloep)"`
	z.Completion(parser.Model.Node, "t1")
	ok := bytes.Contains(z.Out(), []byte(exp))
	if !ok {
		t.Fatalf("expected %s to be present, but did not found it", exp)
	}
}
