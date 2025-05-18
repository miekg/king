# Completion generator for kong

[kong](https://github.com/alecthomas/kong) is a very nice command-line parser for Go. But it misses the
ability to generate (good) shell completions. There are some integrations but they require source level
changes. With _king_ you can just generate the completion files separately from a `kong.Kong`.

This package is copied from [gum](https://github.com/charmbracelet/gum) and made into a standalone library +
some extra features, like telling (via struct tags) how certain things must be completed.

Any struct field can have 2 extra tags:

- `completion:""` which contains a shell command that should be used for completion _or_ a string between
  `<` and `>` which should be a bash action as specified in the complete function in bash(1), like `<file>`
  or `<directory>`.
- `compname:""` not needed often, but this is displayed by (only?) zsh on what type of completion is being
  performed. I.e. it could be "users" when the completion is trying to find users. If not given it the name of
  kong.Node is used (usually the field name in the struct).

I use [zsh](https://zsh.org), so this is where my initial focus is. The
[bash](https://www.gnu.org/software/bash/) completion works, but can probably be done a lot better.

## TODO

- Global options, have a way to add them.
- Bash options, env options also for zsh

## Supported "actions"

The following actions are supported:

- "file", "directory"
- "group"
- "user"
- "export"

And are converted to the correct construct in the completion that is generated.
