package king

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// commandName returns the full path of the kong node. Any alias is ignored.
func commandName(n *kong.Node) (out string) {
	root := n
	for root.Parent != nil {
		root = root.Parent
	}
	return root.Name + identifier(n)
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

// completion returns the completion.
func completion(cmd *kong.Value) string {
	comp := cmd.Tag.Get("completion")
	if comp == "" {
		return ""
	}
	if strings.HasPrefix(comp, "<") && strings.HasSuffix(comp, ">") {
		switch comp {
		case "file", "directory":
			return "_files"
		case "group":
			return "_groups"
		case "user":
			return "_users"
		case "export":
			return "_parameters"
		}
	}
	return comp
}

// compname returns the compname, if there is no compname, cmd.Name is returned.
func compname(cmd *kong.Value) string {
	compname := cmd.Tag.Get("compname")
	if compname == "" {
		return cmd.Name
	}
	return compname
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
