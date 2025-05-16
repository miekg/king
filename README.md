# Completion generator for kong

[kong](https://github.com/alecthomas/kong) is a very nice command-line parser for Go. But it misses the
ability to generate (good) shell completions. There are some integrations but they require source level
changes. With _king_ you can just generate the completion files separately from a `kong.Kong`.

This package is copied from [gum](https://github.com/charmbracelet/gum) and made into a standalone library +
some extra features, like telling (via struct tags) how certain things must be completed. Also see godoc.

I use [zsh](https://zsh.org), so this is where my initial focus is.
