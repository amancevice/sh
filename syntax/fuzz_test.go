//go:build dev.fuzz
// +build dev.fuzz

package syntax_test

import (
	"os/exec"
	"testing"

	"mvdan.cc/sh/v3/syntax"
)

func FuzzQuote(f *testing.F) {
	// Keep in sync with ExampleQuote.
	f.Add("foo")
	f.Add("bar $baz")
	f.Add(`"won't"`)
	f.Add(`~/home`)
	f.Add("#1304")
	f.Add("name=value")
	f.Add(`glob-*`)
	f.Add("invalid-\xe2'")
	f.Add("nonprint-\x0b\x1b")
	f.Fuzz(func(t *testing.T, s string) {
		quoted, ok := syntax.Quote(s)
		if !ok {
			// Contains a null byte; not interesting.
			return
		}
		out, err := exec.Command("bash", "-c",
			"printf %s "+quoted+
				"", // TODO: quoted as a command (func), not an arg
		).CombinedOutput()
		if err != nil {
			t.Fatalf("bash error on %q quoted as %s: %v: %s", s, quoted, err, out)
		}
		want, got := s, string(out)
		if want != got {
			t.Fatalf("output mismatch on %q quoted as %s: got %q (len=%d)", want, quoted, got, len(got))
		}
	})
}
