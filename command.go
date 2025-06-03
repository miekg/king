// Package king is used to generate completions for https://github.com/alecthomas/kong.
// Unlike most other completers this package also completes positional arguments for both Bash and Zsh.
// It can also be used to generate manual pages for a kong command line.
package king

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/kong"
)

// The Completer interface must be implemented by every shell completer. It mainly serves for documentation.
type Completer interface {
	// Completion generates the completion for a shell starting with k. The altname - if not empty - takes
	// precedence over k.Name.
	Completion(k *kong.Node, altname string)
	// Out returns the generated shell completion script.
	Out() []byte
	// Write writes the generated shell completion script to the appropiate file, for Zsh this is _exename and
	// for Bash this is exename.bash, for manual this is m.name.m.Section. If the optional writer is given it
	// contents is written to that.
	Write(w ...io.Writer) error
}

var (
	_ Completer = (*Zsh)(nil)
	_ Completer = (*Bash)(nil)
)

// commandName returns the name of the command node, it takes name from the cmd tag, if that is empty the
// node's name is returned.
func commandName(n *kong.Node) string {
	name := n.Tag.Get("cmd")
	if name != "" {
		return name
	}
	return n.Name
}

func nodePath(n *kong.Node, path []string) *kong.Node {
	if len(path) == 0 {
		return n
	}

	for _, c := range n.Children {
		if c.Name == path[0] {
			if x := nodePath(c, path[1:]); x != nil {
				return x
			}
		}
	}
	return nil
}

// funcName returns the full path of the kong node for use as a function name. Any alias is ignored.
func funcName(n *kong.Node) (out string) {
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
		comp := comp[1 : len(comp)-1]
		return toAction(comp, shell)
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

// toAction returns the proper action per shell.
func toAction(action, shell string) string {
	switch shell {
	case "zsh":
		return zshActions[action]
	case "bash":
		return action
	case "fish":
		return ""
	}
	return ""
}

var zshActions = map[string]string{
	"file":      "_files",
	"directory": "_files",
	"group":     "_groups",
	"user":      "_users",
	"export":    "_parameters",
}
