package king

// Option is an option that can be given to Bash, Zsh or Man.
type Option interface {
	Apply(string) string
}

// RemoveMinus removes any in-command "-".
type RemoveMinus struct{}

func (r RemoveMinus) Apply(s string) string {
}
