package king

import (
	"bytes"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/alecthomas/kong"
)

// Man is a manual page generator.
type Man struct {
	name      string
	Section   int // See mmark's documentation
	Area      string
	WorkGroup string
	Template  string // If empty [ManTemplate] is used.
	manual    []byte
	Flags     []*kong.Flag // Any global flags that the should Application Node have. There are documented after the normal options.
}

// ManTemplate is the default manual page template used when generating a manual page. Where each function
// used is:
//
//   - name: generate a <cmdname> - <single line synopsis>. The node's help is used for this.
//   - synopsis: shows the argument and options. If the node has aliases they are shown here as well.
//   - description: the description of the command's function. The node (non-Kong) "description" tag is used for
//     this.
//   - arguments: a rundown of each of the commands and/or arguments this command has.
//   - options: a list documenting each of the options.
//   - globals: any global flags, from m.Flags.
const ManTemplate = `{{name -}}

{{synopsis -}}

{{description -}}

{{arguments -}}

{{options -}}

{{globals -}}
`

// Out returns the manual in markdown form.
func (m *Man) Out() []byte { return m.manual }

// Write writes the manual page in man format to the file m.name.m.section.
func (m *Man) Write() error {
	if m.manual == nil {
		return fmt.Errorf("no manual")
	}
	// convert to manpage.
	return os.WriteFile(fmt.Sprintf("%s.%d", m.name, m.Section), m.manual, 0644)
}

// nam implements the template func name.
func nam(k, cmd *kong.Node) string {
	help := strings.TrimSuffix(cmd.Help, ".")
	return fmt.Sprintf("## Name\n\n%s - %s\n\n", cmd.Tag.Get("cmd"), help)
}

func synopsis(k, cmd *kong.Node) string {
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
		if f.Required {
			optstring += " " + f.Name
			if f.PlaceHolder != "" {
				optstring += " *" + strings.ToUpper(f.PlaceHolder) + "*"
			}
		}
	}
	fmt.Fprintf(s, "## Synopsis\n\n")
	cmdname := k.Name
	aliases := strings.Join(cmd.Aliases, "|")
	if aliases != "" {
		aliases += "|"
	}
	fmt.Fprintf(s, "`%s`%s%s\n\n", cmdname, optstring, argstring)
	if cmd.Aliases != nil {
		for _, alias := range cmd.Aliases {
			if alias == cmdname {
				continue // already done
			}
			fmt.Fprintf(s, "`%s`%s%s\n\n", alias, optstring, argstring)
		}
	}
	fmt.Fprintf(s, "`%s` %s%s%s%s\n\n", "c", aliases, cmdname, optstring, argstring)
	return s.String()
}

func arguments(k, cmd *kong.Node) string {
	s := &strings.Builder{}
	for _, p := range cmd.Positional {
		formatArg(s, p)
	}
	return s.String()
}

// options implements the options func name.
func options(k, cmd *kong.Node) string {
	s := &strings.Builder{}
	flags := cmd.Flags

	if len(flags) > 0 {
		sort.Slice(flags, func(i, j int) bool { return flags[i].Name < flags[j].Name })
		fmt.Fprintf(s, "### Options\n\n")

		// groups holds any grouped options
		groups := map[string][]*kong.Flag{}
		for _, f := range flags {
			if f.Group != nil {
				groups[f.Group.Key] = append(groups[f.Group.Key], f)
			}
		}

		for _, f := range flags {
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
			fmt.Fprintf(s, "#### %s OPTIONS\n", strings.ToUpper(group))
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

// Manual generates a manual page for field named field of the node.
// On the node k the following tags are used:
//
//   - cmd:"....": command name, overrides k.<named>.Name
//   - aliases:"...": any extra names that this command has.
//   - help:"...": line used in the NAME section: "cmd" - "help" text, as in "ls - list directory contents" if
//     this text ends in a dot it is removed.
//   - description:".....": The entire description paragraph.
//
// Note that any of these may contain markdown markup. The node k doesn't need any special tags.
func (m *Man) Manual(k *kong.Node, name, field string) { // add a field?
	var cmd *kong.Node
	for _, c := range k.Children {
		if c.Name == field {
			cmd = c
			break
		}
	}
	k.Name = name
	m.name = name

	if cmd == nil {
		log.Printf("Failed to generate manual page: %q not found as child", field)
		return
	}

	funcMap := template.FuncMap{
		"name":        func() string { return nam(k, cmd) },
		"description": func() string { return cmd.Tag.Get("description") },
		"synopsis":    func() string { return synopsis(k, cmd) },
		"arguments":   func() string { return arguments(k, cmd) },
		// commands?
		"options": func() string { return options(k, cmd) },
		"globals": func() string { return globals(m.Flags) },
	}

	if m.Template == "" {
		m.Template = ManTemplate
	}

	tmpl := &template.Template{}
	var err error

	tmpl = template.New("generated").Funcs(funcMap)
	tmpl, err = tmpl.Parse(m.Template)
	if err != nil {
		log.Printf("Failed to generate manual page: %s", err)
		return
	}

	format := `%%%%%%
title = "%s %d"
area = "%s"
workgroup = "%s"
# generated by king completion (https://github.com/miekg/king) for kong
%%%%%%

`
	b := &bytes.Buffer{}
	fmt.Fprintf(b, format, cmd.Tag.Get("cmd"), m.Section, m.Area, m.WorkGroup)
	if err = tmpl.Execute(b, nil); err != nil {
		log.Printf("Failed to generate manual page: %s", err)
		return
	}
	m.manual = b.Bytes()
}
