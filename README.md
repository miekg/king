# Completion generator for kong

[kong](https://github.com/alecthomas/kong) is a very nice command-line parser for Go. But it misses the
ability to generate (good) shell completions. There are some integrations but they require source level
changes. With _king_ you can just generate the completion files separately from a `kong.Node`.

This package copies from [gum](https://github.com/charmbracelet/gum) and made into a standalone library +
some extra features, like telling (via struct tags) how certain things must be completed. Also the Bash
completions are completely reworked and for both Zsh and Bash positional argument completion also works.

Any struct field can have an extra tag:

- `completion:""` which contains a shell command that should be used for completion _or_ a string between
  `<` and `>` which should be a Bash action as specified in the `complete` function in bash(1), like `<file>`
  or `<directory>`. These are translated to things Zsh understands.

I use [Zsh](https://zsh.org), so this is where my initial focus is. The
[Bash](https://www.gnu.org/software/bash/) completion works, but can probably be done a lot better.

## TODO

- Bash short options and aliases.
- Env flag from Kong

## Supported "actions"

The following actions are supported:

- "file", "directory"
- "group"
- "user"
- "export"

And are converted to the correct construct in the completion that is generated.
