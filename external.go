package displaywidth

import "sync"

var (
	externalMu     sync.RWMutex
	externalWidths map[string]int
)

// SetExternalWidths installs a map of grapheme-cluster-string → cell-width
// that overrides the library's built-in width tables. This is for callers
// that have probed the terminal directly (e.g. via CSI 6n cursor position
// reports) and want String/Bytes/Rune — and consequently downstream
// libraries like charmbracelet/x/ansi and charm.land/lipgloss — to use
// the probed widths instead of the spec-derived ones.
//
// Lookup precedence per grapheme cluster:
//  1. The external map, if non-nil and the cluster is present.
//  2. The library's trie + VS16/skin-tone heuristics.
//
// The map is treated as immutable post-install. Callers must replace
// the map by calling SetExternalWidths again with a new instance rather
// than mutating the existing one. Pass nil to clear.
func SetExternalWidths(m map[string]int) {
	externalMu.Lock()
	defer externalMu.Unlock()
	externalWidths = m
}

// externalLookup returns (width, true) if s is in the external widths
// map, or (0, false) otherwise. Uses a snapshot pattern so the lock is
// held only briefly.
func externalLookup[T ~string | []byte](s T) (int, bool) {
	externalMu.RLock()
	m := externalWidths
	externalMu.RUnlock()
	if m == nil {
		return 0, false
	}
	w, ok := m[string(s)]
	return w, ok
}
