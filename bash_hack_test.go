package king

import "testing"

func TestBashHack(t *testing.T) {
	const in = `'protogrp ch'*'--admin-domains')
      while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -W "$(_c_filter "$(c admindomain --comp)")" -- "$cur")
      ;;
    'protogrp ch'*'-d')
      while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -W "$(_c_filter "$(c admindomain --comp)")" -- "$cur")
      ;;
    'protogrp rm'*)
      while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen -W "$(_c_filter "$(c protogrp list --comp)")" -- "$cur")
      ;;
`
	sorted := sortcmd(in)
	t.Log(sorted)
}
