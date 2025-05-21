package king

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/kong"
)

// formatFlag is used to format an option. If quote is given and true the whole thing is indented.
func formatFlag(s io.Writer, f *kong.Flag, quote ...bool) {
	q := ""
	if len(quote) > 0 {
		q = "> "
	}
	if f.Tag.Has("negatable") {
		fmt.Fprintf(s, "%s`--[no-]%s`", q, f.Name)
	} else {
		fmt.Fprintf(s, "%s`--%s`", q, f.Name)
	}
	if f.Short != 0 {
		fmt.Fprintf(s, ", `-%c`", f.Short)
	}
	// TODO(miek): probably even more types then are now here.
	// TODO(miek): entire rework this; what does it even attempt? And how can we make this do the right thing
	// for custom types.
	switch f.Target.Kind().String() {
	case "Size", "slice", "string", "int", "uint", "uint8", "uint16", "uint32", "uint64":
		// if the type:counter is set, this option can be repeated, and cannot have a value
		if f.Tag.Get("type") == "counter" {
			break
		}
		if f.PlaceHolder != "" {
			fmt.Fprintf(s, " *%s*", strings.ToUpper(f.PlaceHolder))
			break
		}
		fmt.Fprintf(s, " *%s*", strings.TrimPrefix(strings.ToUpper(f.Name), "--"))
	case "Time":
		if f.Format != "" {
			// convert time Format modifiers back into something more sane.
			f.Format = strings.Replace(f.Format, "2006", "yyyy", 1)
			f.Format = strings.Replace(f.Format, "01", "mm", 1)
			f.Format = strings.Replace(f.Format, "02", "dd", 1)
			f.Format = strings.Replace(f.Format, "15", "hh", 1)
			f.Format = strings.Replace(f.Format, "04", "mm", 1) // clashes with month, but we uppercase this ....

			fmt.Fprintf(s, " *%s*", strings.ToUpper(f.Format))
			break
		}
		fmt.Fprintf(s, " *%s*", strings.TrimPrefix(strings.ToUpper(f.Name), "--"))
	case "Date":
		fmt.Fprintf(s, " *%s*", "YYYY-MM-DD|[+-]Y")
	}
	fmt.Fprintln(s)
	deprecated := ""
	if f.Tag.Has("deprecated") {
		deprecated = "(Deprecated) "
	}
	fmt.Fprintf(s, "%s:   %s%s", q, deprecated, f.Help)
	if f.Required {
		fmt.Fprintf(s, " This is a required option.")
	}
	if f.Tag.Get("type") == "counter" {
		fmt.Fprintf(s, " This option can be repeated.")
	}
	if f.Enum != "" {
		enums := f.EnumSlice()
		switch len(enums) {
		case 1:
			fmt.Fprintf(s, " Valid value is: ")
			fmt.Fprintf(s, "%q.", enums[0])
		case 2:
			fmt.Fprintf(s, " Valid values are: ")
			fmt.Fprintf(s, "%q or %q.", enums[0], enums[1])
		default:
			fmt.Fprintf(s, " Valid values are: ")
			div := ", "
			for i, e := range enums {
				if i == len(enums)-2 {
					div = " or "
				}
				if i == len(enums)-1 {
					div = "."
				}
				fmt.Fprintf(s, "%q%s", e, div)
			}
		}
		if f.Default != "" {
			fmt.Fprintf(s, " The default is %q.", f.Default)
		}
	}
	if f.Enum == "" && f.Default != "" { // No enum, but still a default
		fmt.Fprintf(s, " The default is %q.", f.Default)
	}

	if f.Envs != nil {
		for i := range f.Envs {
			f.Envs[i] = "`${" + f.Envs[i] + "}`"
		}
		vars := "variables"
		if len(f.Envs) == 1 {
			vars = "variable"
		}
		fmt.Fprintf(s, " The default value is derived from the environment %s: %s.", vars, strings.Join(f.Envs, ", "))
	}

	if f.Xor != nil {
		fmt.Fprintf(s, " This option can not be used together with: ")
		if len(f.Xor) == 2 {
			// div can be "."
			for i := range f.Xor {
				if "--"+f.Xor[i] != f.Name {
					fmt.Fprintf(s, "**--%s**%s", f.Xor[i], ".")
				}
			}
		} else {
			div := ", "
			for i := range f.Xor {
				if "--"+f.Xor[i] != f.Name {
					fmt.Fprintf(s, "**--%s**%s", f.Xor[i], div)
				}
				if i == len(f.Xor)-2 {
					div = "."
				}
			}
		}
	}
	// We allow one alias for an option.
	if len(f.Aliases) > 0 && f.Aliases[0] != "" {
		fmt.Fprintf(s, " This option has **--%s** as an alias.", f.Aliases[0])
	}

	fmt.Fprintln(s)
	fmt.Fprintln(s)
}

func formatArg(s io.Writer, p *kong.Positional) {
	name := p.Name
	if p.Tag.PlaceHolder != "" {
		name = p.Tag.PlaceHolder
	}
	// required opts....
	fmt.Fprintf(s, "`%s`", strings.ToUpper(name))
	if p.Tag.Short != 0 {
		fmt.Fprintf(s, ", `-%c`", p.Tag.Short)
	}
	fmt.Fprintln(s)
	fmt.Fprintf(s, ":   %s", p.Help)
	if !p.Required {
		fmt.Fprintf(s, " This argument is optional.")
	}
	if p.Enum != "" {
		fmt.Fprintf(s, " Valid values are: ")
		enums := p.EnumSlice()
		switch len(p.EnumSlice()) {
		case 1:
			fmt.Fprintf(s, "%q.", enums[0])
		case 2:
			fmt.Fprintf(s, "%q or %q.", enums[0], enums[1])
		default:
			div := ", "
			for i, e := range enums {
				if i == len(enums)-2 {
					div = " or "
				}
				if i == len(enums)-1 {
					div = "."
				}
				fmt.Fprintf(s, "%q%s", e, div)
			}
		}
	}
	if p.Default != "" {
		fmt.Fprintf(s, " The default is %q.", p.Default)
	}

	fmt.Fprintln(s)
}

func formatCmd(s io.Writer, cmd *kong.Command) {
	fmt.Fprintf(s, "%+v\n", cmd)
}
