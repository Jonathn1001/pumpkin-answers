package slug_test

import (
	"strings"
	"testing"

	"claimsplatform/internal/slug"
)

func TestMake(t *testing.T) {
	cases := []struct{ in, want string }{
		{"SafeGuard Insurance", "safeguard-insurance"},
		{"HealthFirst", "healthfirst"},
		{"Bảo Việt Hà Nội", "bao-viet-ha-noi"}, // Vietnamese diacritics stripped
		{"Đặng Trần", "dang-tran"},             // đ/Đ do not decompose under NFD
		{"A & B  --  C!", "a-b-c"},             // punctuation runs collapse to one hyphen
		{"-Hello-", "hello"},                   // leading/trailing hyphens trimmed
		{"123 Main", "123-main"},               // may start with a digit
		{"@@@", ""},                            // nothing usable
		{"", ""},
	}
	for _, c := range cases {
		if got := slug.Make(c.in); got != c.want {
			t.Errorf("Make(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMakeCapsAt63(t *testing.T) {
	got := slug.Make(strings.Repeat("x", 80))
	if len(got) != 63 {
		t.Fatalf("len = %d, want 63", len(got))
	}
}

func TestMakeNeverEndsWithHyphenAfterCap(t *testing.T) {
	// 62 'a', then a space at position 63 would otherwise leave a trailing hyphen.
	got := slug.Make(strings.Repeat("a", 62) + " bbbb")
	if strings.HasSuffix(got, "-") {
		t.Fatalf("slug %q must not end with a hyphen", got)
	}
}
