package king

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/v2/mparser"
	"github.com/mmarkdown/mmark/v2/render/man"
)

// Man is a manual page generator.
type Man struct {
	name      string
	manual    []byte
	Section   int // See mmark's documentation
	Area      string
	WorkGroup string
	Template  string       // If empty [ManTemplate] is used.
	Flags     []*kong.Flag // Any global flags that the should Application Node have. There are documented after the normal flags.
}

// ManTemplate is the default manual page template used when generating a manual page. Where each function
// used is:
//
//   - name: generate a <cmdname> - <single line synopsis>. The node's help is used for this. The removes the final dot from the help, and lowercases
//     the first letter if the word isn't all caps.
//     The command's _altname_ is always used here.
//   - synopsis: shows the argument and options. If the node has aliases they are shown here as well.
//   - description: the description of the command's function. The node's (non-Kong) "description" tag is used for
//     this.
//   - commands: a rundown of each of the subcommands this command has.
//   - arguments: a rundown of each of the arguments this command has.
//   - options: a list documenting each of the options.
//   - globals: any global flags, from m.Flags.
//
// Note that the TOML manual header is always used.
const ManTemplate = `{{name -}}

{{synopsis -}}

{{description -}}

{{commands -}}

{{arguments -}}

{{options -}}

{{globals -}}
`

// Out returns the manual in markdown form.
func (m *Man) Out() []byte { return m.manual }

// Write writes the manual page in man format to the file m.name.m.section.
func (m *Man) Write(w ...io.Writer) error {
	if m.manual == nil {
		return fmt.Errorf("no manual")
	}
	p := parser.NewWithExtensions(parser.FencedCode | parser.DefinitionLists | parser.Tables)
	p.Opts = parser.Options{
		ParserHook: func(data []byte) (ast.Node, []byte, int) { return mparser.Hook(data) },
		Flags:      parser.FlagsNone,
	}
	doc := markdown.Parse(m.manual, p)
	renderer := man.NewRenderer(man.RendererOptions{})
	md := markdown.Render(doc, renderer)
	if len(w) > 0 {
		w[0].Write(md)
	}
	return os.WriteFile(fmt.Sprintf("%s.%d", m.name, m.Section), md, 0644)
}

// Manual generates a manual page for child node that can be found via field, where field may contain
// a space seperated list of node names: "mfa list", looks for the mfa node and its child named list.
// On the node k the following tags are used:
//
//   - cmd:"....": command name, overrides k.<named>.Name
//   - aliases:"...": any extra names that this command has.
//   - help:"...": line used in the NAME section: "cmd" - "help" text, as in "ls - list directory contents" if
//     this text ends in a dot it is removed.
//   - description:".....": The entire description paragraph.
//
// Note that any of these may contain markdown markup. The node k doesn't need any special
//
// If path is empty, the manual page for k is returned, altname is used as an alternative name for this command, which
// is not found in any of the tags. This is also the default name under which the manual page is saved.
//
// When path is is not empty it should have a list of command names which are followed to find the node
// for which we generate the manual page. This path is also used to construct the command line(s) in the synopsis.
//
// rootname is the name of the main executable, this can't be put as a tag on the main node, as that node itself can't carry struct tags.
//
// To generate a manual page for the main node, `k` we need a description and help, this can be done via:
//
//	type WrapT struct {
//		Wrap MyCli `cmd:"" help:"my help" description:"my desc"`
//	}
//
// and then call Manual with: m.Manual(parser.Model.Node, "_wrap", "MyExec", "")
// The prevent "wrap" showing up as a valid command, is can be prefixed it with _. This only checked on the
// first element in path.
func (m *Man) Manual(k *kong.Node, path, altname, rootname string) {
	fields := strings.Fields(path)
	if strings.HasPrefix(fields[0], "_") {
		fields[0] = fields[0][1:]
	}
	cmd := nodePath(k, fields)
	if cmd == nil {
		return
	}
	m.name = altname

	funcMap := template.FuncMap{
		"name":        func() string { return name(cmd, altname, rootname) },
		"description": func() string { return description(cmd) },
		"synopsis":    func() string { return m.synopsis(cmd, path, altname, rootname) },
		"arguments":   func() string { return arguments(cmd) },
		"commands":    func() string { return commands(cmd) },
		"options":     func() string { return options(cmd) },
		"globals":     func() string { return globals(m.Flags) },
	}

	if m.Template == "" {
		m.Template = ManTemplate
	}

	tmpl := &template.Template{}
	var err error

	tmpl = template.New("manualpage").Funcs(funcMap)
	tmpl, err = tmpl.Parse(m.Template)
	if err != nil {
		log.Printf("Failed to generate manual page: %s", err)
		return
	}

	format := `%%%%%%
title = "%s %d"
area = "%s"
workgroup = "%s"
# generated by king (https://github.com/miekg/king) for kong
%%%%%%

`
	b := &bytes.Buffer{}
	name := altname
	if name == "" {
		name = fmt.Sprintf("%s %s", rootname, path)
	}
	fmt.Fprintf(b, format, name, m.Section, m.Area, m.WorkGroup)
	if err = tmpl.Execute(b, nil); err != nil {
		log.Printf("Failed to generate manual page: %s", err)
		return
	}
	m.manual = b.Bytes()
}

// name implements the template func name.
func name(cmd *kong.Node, altname, rootname string) string {
	help := strings.TrimSuffix(cmd.Help, ".")
	if strings.ToUpper(help) != help && len(help) > 2 { // not all caps
		help = strings.ToLower(help[0:1]) + help[1:]
	}
	if altname == "" {
		altname = rootname + " " + commandName(cmd)
	}
	return fmt.Sprintf("## Name\n\n%s - %s\n\n", altname, help)
}

func (m *Man) synopsis(cmd *kong.Node, path, altname, rootname string) string {
	s := &strings.Builder{}

	optstring := " *[OPTION]*"
	if len(cmd.Flags) == 0 {
		optstring = ""
	}
	if len(cmd.Flags) > 0 {
		optstring += "..."
	}

	argstring := ""
	for _, a := range cmd.Positional {
		name := a.Name
		if a.Tag.PlaceHolder != "" {
			name = a.Tag.PlaceHolder
		}
		if a.Required {
			argstring += " *" + strings.ToUpper(name) + "*"
		} else {
			argstring += " *[" + strings.ToUpper(name) + "]*"
		}
	}
	for _, f := range cmd.Flags {
		if f.Hidden {
			continue
		}
		if f.Required {
			optstring += " --" + f.Name
			if f.PlaceHolder != "" {
				optstring += " *" + strings.ToUpper(f.FormatPlaceHolder()) + "*"
			}
		}
	}
	cmdstring := ""
	for _, c := range cmd.Children {
		if c.Hidden {
			continue
		}
		if c.Type != kong.CommandNode {
			continue
		}
		cmdname := commandName(c)
		if cmdstring != "" {
			cmdstring += "|"
		} else {
			cmdstring = " "
		}
		cmdstring += cmdname
	}
	if len(cmdstring) > 40 { // dumb check, but we can have a large number of subcommands
		cmdstring = " *[COMMAND]*..."
	}

	// drop the last element in path
	ignore := strings.HasPrefix(path, "_")
	fields := strings.Fields(path)
	path = ""
	if len(fields) > 0 {
		path = strings.Join(fields[:len(fields)-1], " ")
	}
	if rootname != "" {
		rootname += " "
	}
	if path != "" {
		path += " "
	}
	fmt.Fprintf(s, "## Synopsis\n\n")
	if altname != "" {
		fmt.Fprintf(s, "`%s`%s%s%s\n\n", altname, optstring, argstring, cmdstring)
	}
	if !ignore {
		fmt.Fprintf(s, "`%s%s%s`%s%s%s\n\n", rootname, path, commandName(cmd), optstring, argstring, cmdstring)
		for _, alias := range cmd.Aliases {
			fmt.Fprintf(s, "`%s%s%s`%s%s%s\n\n", rootname, path, alias, optstring, argstring, cmdstring)
		}
	}
	return s.String()
}

func description(cmd *kong.Node) string {
	s := &strings.Builder{}
	fmt.Fprint(s, "## Description\n\n")
	fmt.Fprint(s, cmd.Tag.Get("description"))
	fmt.Fprintln(s)
	return s.String()
}

func arguments(cmd *kong.Node) string {
	if !hasPositional(cmd) {
		return ""
	}
	s := &strings.Builder{}
	fmt.Fprintf(s, "\nThe following positional arguments are available:\n\n")
	for _, p := range cmd.Positional {
		// hidden!
		formatArg(s, p)
	}
	return s.String()
}

func commands(cmd *kong.Node) string {
	if !hasCommands(cmd) {
		return ""
	}
	s := &strings.Builder{}
	fmt.Fprintf(s, "\nThe following subcommands are available:\n\n")
	for _, c := range cmd.Children {
		if c.Type == kong.CommandNode && !c.Hidden {
			formatCmd(s, c)
		}
	}
	return s.String()
}

// options implements the options func name.
func options(cmd *kong.Node) string {
	s := &strings.Builder{}
	flags := cmd.Flags

	if len(flags) > 0 {
		sort.Slice(flags, func(i, j int) bool { return flags[i].Name < flags[j].Name })
		fmt.Fprintf(s, "### Options\n\n")

		// groups holds any grouped options
		groups := map[string][]*kong.Flag{}
		for _, f := range flags {
			if f.Hidden {
				continue
			}
			if f.Group != nil {
				groups[f.Group.Key] = append(groups[f.Group.Key], f)
			}
		}

		for _, f := range flags {
			if f.Hidden {
				continue
			}
			if f.Group == nil {
				formatFlag(s, f)
			}
		}
		fmt.Fprintln(s)
		// format groups options
		for k := range groups {
			if k != strings.ToLower(k) {
				log.Fatalf("Group keys must be all lowercase: %s", k)
			}
		}
		keys := slices.Sorted(maps.Keys(groups))

		for _, group := range keys {
			fmt.Fprintf(s, "#### %s Options\n\n", group)
			for _, f := range groups[group] {
				formatFlag(s, f, true)
			}
		}
	}
	return s.String()
}

func globals(flags []*kong.Flag) string {
	s := &strings.Builder{}
	if len(flags) > 0 {
		fmt.Fprintf(s, "The following default options are available.\n\n")
		for _, f := range flags {
			formatFlag(s, f)
		}
	}
	return s.String()
}
