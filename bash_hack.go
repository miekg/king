package king

import (
	"bufio"
	"slices"
	"strings"
)

// sortcmd is a dump hack to sort the bash 'case' cmd contruct from longtest to shortest, to make the matching
// work.
func sortcmd(buf string) string {
	// This buffer consists out of a sequence of 3 lines, there the first line is the thing we need to sort
	// on. After that we just output it in the correct order.
	//
	//'website ch'*)
	// while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -W "$(_c_filter "--force -f --url --status --volume --comment --user --group $(c website list --comp)")" -- "$cur")
	//;;

	scanner := bufio.NewScanner(strings.NewReader(buf))
	cases := map[string]string{}

	i := 0
	cur := ""
	for scanner.Scan() {
		t := scanner.Text()
		switch i % 3 {
		case 0:
			cur = t
			cases[cur] = t + "\n"
		case 1:
			cases[cur] += t + "\n"
		case 2:
			cases[cur] += t + "\n"
		}
		i++
	}
	keys := []string{}
	for k := range cases {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b string) int {
		if x := len(b) - len(a); x != 0 {
			return x
		}
		return strings.Compare(a, b)
	})
	b := &strings.Builder{}
	for _, k := range keys {
		b.WriteString(cases[k])
	}
	return b.String()
}
