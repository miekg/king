// Package king is used to generate completions for https://github.com/alecthomas/kong.
// Unlike most other completers this package also completes positional arguments for both Bash and Zsh.
package king

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/kong"
)

// The Completer interface must be implemented by every shell completer. It mainly serves for documentation.
type Completer interface {
	// Completion generates the completion for a shell starting with k. The exename - if not empty is the name of the executable, if it
	// is different from k.Name.
	Completion(k *kong.Node, exename string)
	// Out returns the generated shell completion script.
	Out() []byte
	// Write writes the generated shell completion script to the appropiate file, for Zsh this is _exename and for Bash this is exename.bash
	Write() error
}

var (
	_ Completer = (*Zsh)(nil)
	_ Completer = (*Bash)(nil)
)

// commandName returns the full path of the kong node. Any alias is ignored.
func commandName(n *kong.Node) (out string) {
	root := n
	for root.Parent != nil {
		root = root.Parent
	}
	out = strings.Replace(root.Name+identifier(n), ".", "_", -1)
	return strings.Replace(out, "-", "_", -1)
}

// identifier creates a name suitable for using as an identifier in shell code.
func identifier(n *kong.Node) (out string) {
	if n.Parent != nil {
		out += identifier(n.Parent)
	}
	switch n.Type {
	case kong.CommandNode:
		out += "_" + n.Name
	case kong.ArgumentNode:
		out += "_" + n.Name
	default:
	}
	return out
}

func hasCommands(cmd *kong.Node) bool {
	for _, c := range cmd.Children {
		if !c.Hidden {
			return true
		}
	}
	return false
}

// hasPositional returns true if there are positional arguments.
func hasPositional(cmd *kong.Node) bool { return len(cmd.Positional) > 0 }

// completions returns all completions that this kong.Node has.
func completions(cmd *kong.Node) []string {
	completions := []string{}
	for _, c := range cmd.Children {
		if c.Hidden {
			continue
		}
		completions = append(completions, c.Name)
	}
	for _, f := range cmd.Flags {
		if f.Hidden {
			continue
		}
		completions = append(completions, "--"+f.Name)
		if f.Short != 0 {
			completions = append(completions, "-"+fmt.Sprintf("%c", f.Short))
		}
	}
	for _, p := range cmd.Positional {
		completions = append(completions, completion(p, "bash"))
	}
	return completions
}

// completion returns the completion for the shell for the kong.Value.
func completion(cmd *kong.Value, shell string) string {
	comp := cmd.Tag.Get("completion")
	if comp == "" {
		return ""
	}
	if strings.HasPrefix(comp, "<") && strings.HasSuffix(comp, ">") {
		comp := comp[1 : len(comp)-2]
		switch comp {
		case "file", "directory":
			if shell == "zsh" {
				return "_files"
			}
			return comp
		case "group":
			if shell == "zsh" {
				return "_groups"
			}
			return comp
		case "user":
			if shell == "zsh" {
				return "_users"
			}
			return comp
		case "export":
			if shell == "zsh" {
				return "_parameters"
			}
			return comp
		}
	}
	return "$(" + comp + ")"
}

// writeString writes a string into a buffer, and checks if the error is not nil.
func writeString(b io.StringWriter, s string) { b.WriteString(s) }

func flagEnums(flag *kong.Flag) []string {
	values := make([]string, 0)
	for _, enum := range flag.EnumSlice() {
		if strings.TrimSpace(enum) != "" {
			values = append(values, enum)
		}
	}
	return values
}

func flagEnvs(flag *kong.Flag) []string {
	values := make([]string, 0)
	for _, env := range flag.Envs {
		if strings.TrimSpace(env) != "" {
			values = append(values, env)
		}
	}
	return values
}
