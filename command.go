package king

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

type Completer interface {
	Completion(*kong.Node, string)
	Out() []byte
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

// hasPositional returns true if there are positional arguments.
func hasPositional(cmd *kong.Node) bool { return len(cmd.Positional) > 0 }

// completion returns the completion for the shell.
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
func writeString(b io.StringWriter, s string) {
	if _, err := b.WriteString(s); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func flagEnums(flag *kong.Flag) []string {
	values := make([]string, 0)
	for _, enum := range flag.EnumSlice() {
		if strings.TrimSpace(enum) != "" {
			values = append(values, enum)
		}
	}
	return values
}
