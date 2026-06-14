// Package comparison turns two config documents into a field-level changelog,
// wrapping r3labs/diff (library-first; no hand-rolled diff).
package comparison

import (
	"strings"

	"claimsplatform/internal/domain"

	"github.com/r3labs/diff/v3"
)

type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Change is the display-friendly diff entry the API returns.
type Change struct {
	Path  string     `json:"path"`
	Type  ChangeType `json:"type"`
	Left  any        `json:"left,omitempty"`
	Right any        `json:"right,omitempty"`
}

func Diff(left, right domain.ConfigDocument) ([]Change, error) {
	changelog, err := diff.Diff(left, right, diff.SliceOrdering(true), diff.TagName("json"))
	if err != nil {
		return nil, err
	}
	out := make([]Change, 0, len(changelog))
	for _, c := range changelog {
		out = append(out, Change{
			Path:  strings.Join(c.Path, "."),
			Type:  mapType(c.Type),
			Left:  c.From,
			Right: c.To,
		})
	}
	return out, nil
}

func mapType(t string) ChangeType {
	switch t {
	case diff.CREATE:
		return Added
	case diff.DELETE:
		return Removed
	default:
		return Changed
	}
}
