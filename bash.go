package king

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// Bash is a bash completion generator.
type Bash struct {
	name       string
	completion []byte
	Flags      []*kong.Flag // Any global flags that the should Application Node have.
}

func (b *Bash) Out() []byte { return b.completion }

func (b *Bash) Write(w ...io.Writer) error {
	if b.completion == nil {
		return fmt.Errorf("no completion")
	}
	if len(w) > 0 {
		w[0].Write(b.completion)
	}
	return os.WriteFile(b.name+".bash", b.completion, 0644)
}

func (b *Bash) Completion(k *kong.Node, altname string) {
	k.Flags = append(k.Flags, b.Flags...)
	format := `# bash completion for %[1]s
# generated by king (https://github.com/miekg/king) for kong

`
	var out strings.Builder
	if altname == "" {
		b.name = k.Name
	} else {
		b.name = altname
		k.Name = altname
	}
	fmt.Fprintf(&out, format, b.name)
	b.gen(&out, k)
	b.completion = []byte(out.String())
}

func (b Bash) writeFilterFunc(buf io.StringWriter) {
	format := `_%[1]s_filter() {
  COMP_REPLY=()
  local words="$1"
  local cur=${COMP_WORDS[COMP_CWORD]}
  local result=()

  if [[ "${cur:0:1}" == "-" ]]; then
    echo "$words"

  else
    for word in $words; do
      [[ "${word:0:1}" != "-" ]] && result+=("$word")
    done
    echo "${result[*]}"
  fi
}
`
	writeString(buf, fmt.Sprintf(format, b.name))
}

func (b Bash) compReply(completions []string) string {
	if len(completions) == 1 && !strings.HasPrefix(completions[0], "$") && !strings.HasPrefix(completions[0], "--") { // action and not empty
		format := `while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -A %s -- "$cur")` + "\n"
		return fmt.Sprintf(format, completions[0])
	}
	format := `while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -W "$(_%s_filter "%s")" -- "$cur")` + "\n"
	return fmt.Sprintf(format, b.name, strings.Join(completions, " "))
}

func (b Bash) writeFlag(buf io.StringWriter, flag *kong.Flag, parents ...string) {
	if flag.Hidden {
		return
	}
	p := ""
	if len(parents) > 0 {
		p = parents[0] + " "
	}
	// 'user add'*'--backup-backend')
	//   while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -W "$(_xxx_filter "s3")" -- "$cur")
	//   ;;
	completions := []string{}
	if enums := flagEnums(flag); len(enums) > 0 {
		completions = enums
	}
	if comptag := completion(flag.Value, "bash"); comptag != "" {
		completions = []string{comptag}
	}
	if envs := flagEnvs(flag); len(envs) > 0 {
		for i := range envs {
			envs[i] = "$" + envs[i]
		}
		completions = envs
	}

	if len(completions) == 0 { // nothing to complete
		return
	}
	writeString(buf, fmt.Sprintf(`    '%s'*'--%s')`+"\n", strings.TrimSpace(p), flag.Name))
	writeString(buf, "      "+b.compReply(completions))
	writeString(buf, "      ;;\n")
	if flag.Short != 0 {
		writeString(buf, fmt.Sprintf(`    '%s'*'-%c')`+"\n", strings.TrimSpace(p), flag.Short))
		writeString(buf, "      "+b.compReply(completions))
		writeString(buf, "      ;;\n")
	}
}

// writeCommand writes a completion case statement. The optional parent is used to create the correct matching
// for sub-sub comments
func (b Bash) writeCommand(buf io.StringWriter, cmd *kong.Node, parents ...string) {
	p := ""
	if len(parents) > 0 {
		p = parents[0] + " "
	}
	if cmd.Type == kong.ApplicationNode {
		for _, c := range cmd.Children {
			b.writeCommand(buf, c, p+c.Name)
		}
		return
	}
	//'group add'*)
	//  while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -W "$(_xxx_filter "--gid --auto --man --help")" -- "$cur")
	//  ;;
	writeString(buf, fmt.Sprintf(`    '%s'*)`+"\n", strings.TrimSpace(p)))
	completions := completions(cmd)
	writeString(buf, "      "+b.compReply(completions))
	writeString(buf, "      ;;\n")
	for _, f := range cmd.Flags {
		b.writeFlag(buf, f, p)
	}
	for _, c := range cmd.Children {
		b.writeCommand(buf, c, p+c.Name)
	}
}

func (b Bash) writeApp(buf io.StringWriter, cmd *kong.Node) {
	completions := completions(cmd)
	writeString(buf, "      "+b.compReply(completions))
	writeString(buf, "      ;;\n")
}

func (b Bash) gen(buf io.StringWriter, cmd *kong.Node) {
	b.writeFilterFunc(buf)

	cmdName := funcName(cmd)
	if b.name != "" {
		writeString(buf, fmt.Sprintf("\n_%s_completions() {\n", b.name))
	} else {
		writeString(buf, fmt.Sprintf("\n_%s_completions() {\n", cmdName))
	}
	writeString(buf, `  local cur=${COMP_WORDS[COMP_CWORD]}
  local compwords=("${COMP_WORDS[@]:1:$COMP_CWORD-1}")
  local compline="${compwords[*]}"

  case "$compline" in

`)
	if hasCommands(cmd) {
		newbuf := &strings.Builder{}
		b.writeCommand(newbuf, cmd)
		sorted := sortcmd(newbuf.String())
		writeString(buf, sorted)
	}
	for _, f := range cmd.Flags {
		b.writeFlag(buf, f)
	}
	writeString(buf, "\n"+`    *)`+"\n")

	if hasPositional(cmd) && len(cmd.Positional) > 1 {
		writeString(buf, "      "+`COMP_CARG=$COMP_CWORD; for i in "${COMP_WORDS[@]}"; do [[ ${i} == -* ]] && ((COMP_CARG = COMP_CARG - 1)); done`+"\n")
		writeString(buf, "      "+`case $COMP_CARG in`+"\n")

		for i, p := range cmd.Positional {
			writeString(buf, fmt.Sprintf("\n"+`        '%d')`+"\n", i+1))
			comptag := []string{completion(p, "bash")}
			writeString(buf, "          "+b.compReply(comptag))
			writeString(buf, "          return\n          ;;\n")
		}

		writeString(buf, "\n      "+`esac`+"\n\n")
	}

	b.writeApp(buf, cmd)

	writeString(buf, `esac

} &&`)
	if b.name != "" {
		writeString(buf, fmt.Sprintf("\ncomplete -F _%[1]s_completions %[1]s\n", b.name))
	} else {
		writeString(buf, fmt.Sprintf("\ncomplete -F _%[1]s_completions %[1]s\n", cmdName))
	}
}
