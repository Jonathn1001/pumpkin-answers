// Package slug derives URL-safe tenant slugs from human display names. Slug
// generation is a domain concern: it lives here so every client (admin UI,
// direct API, onboarding scripts) gets identical, server-authoritative slugs.
package slug

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const maxLen = 63 // matches backend slug rule ^[a-z0-9][a-z0-9-]{0,62}$

// Make converts a display name into a slug matching ^[a-z0-9][a-z0-9-]{0,62}$.
// Diacritics (including Vietnamese) are stripped, runs of non-alphanumerics
// collapse to a single hyphen, and the result is capped at 63 chars with no
// leading/trailing hyphen. Returns "" when nothing usable remains.
func Make(name string) string {
	// đ/Đ are distinct letters that do not decompose under NFD; map them first.
	name = strings.NewReplacer("đ", "d", "Đ", "d").Replace(name)
	var b strings.Builder
	pendingHyphen := false
	for _, r := range norm.NFD.String(name) {
		if unicode.Is(unicode.Mn, r) { // combining diacritical mark
			continue
		}
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'z':
			if pendingHyphen && b.Len() > 0 {
				b.WriteByte('-')
			}
			pendingHyphen = false
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			if pendingHyphen && b.Len() > 0 {
				b.WriteByte('-')
			}
			pendingHyphen = false
			b.WriteRune(r - 'A' + 'a')
		default:
			pendingHyphen = true // defer: avoids leading/trailing/double hyphens
		}
		if b.Len() >= maxLen {
			break
		}
	}
	s := b.String()
	if len(s) > maxLen {
		s = s[:maxLen]
	}
	return strings.TrimRight(s, "-")
}
