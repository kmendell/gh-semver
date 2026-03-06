package semver

import "testing"

func TestParseSemVerRejectsInvalidTag(t *testing.T) {
	if _, err := ParseSemVer("not-a-semver-tag"); err == nil {
		t.Fatal("expected ParseSemVer to reject invalid tags")
	}
}