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
	return strings.Replace(root.Name+identifier(n), ".", "_", -1)
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

/*
// commands returns all possible paths through the (sub)command structure of the kong.Node.
func commands(node *kong.Node) []string {
	if node == nil {
		return nil
	}

	var paths [][]string
	var dfs func(node *kong.Node, currentPath []string)
	dfs = func(node *kong.Node, currentPath []string) {
		if node.Hidden {
			return
		}

		if node.Type == kong.CommandNode {
			currentPath = append(currentPath, node.Name)
		}

		if len(node.Children) == 0 {
			paths = append(paths, currentPath)
			return
		}

		for _, c := range node.Children {
			dfs(c, currentPath)
		}
	}

	dfs(node, []string{})
	// Now we have all paths, e.g. [[do] [more] [even-more do-even-more] [even-more what-even-more]]
	// But we need all posible ones, e.g. 'do', 'more', 'even-more', 'even-more do-even-more', etc.
	// So iterator over them and split each and add them to a map by increasing the number of fields.
	m := map[string]struct{}{}
	for _, p := range paths {
		f1 := ""
		sep := ""
		for _, p1 := range p {
			fields := strings.Fields(p1)
			for _, f := range fields {
				f1 += sep + f
				m[f1] = struct{}{}
				sep = " "
			}
		}
	}
	keys := slices.Collect(maps.Keys(m))
	slices.SortFunc(keys, func(a, b string) int { return len(b) - len(a) })
	return keys
}

// nodeForCommand uses the trace string (which is a possible command line invocation to find the last node and returns that node.
func nodeForCommand(cmd *kong.Node, trace string) *kong.Node {
	var dfs func(*kong.Node, []string) *kong.Node
	dfs = func(node *kong.Node, fields []string) *kong.Node {
		if node.Name == fields[0] && len(fields) == 1 {
			return node
		}

		if node.Name == fields[0] {
			fields = fields[1:]
		}
		for _, c := range node.Children {
			if n := dfs(c, fields); n != nil {
				return n
			}
		}
		return nil
	}

	fields := strings.Fields(trace)
	leaf := dfs(cmd, fields)
	return leaf
}
*/

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
