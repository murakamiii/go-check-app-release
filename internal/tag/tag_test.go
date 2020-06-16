package tag

import (
	"testing"
)

func TestNewerThanRight(t *testing.T) {
	cases := []struct {
		lhs string
		rhs string
		expects bool
	}{
		{ "0.0.1", "0.0.2", false },
		{ "0.0.3", "0.0.2", true },
		{ "1.1.1", "0.99.99", true },
	}
	for _, c := range cases {
		if newerThanRight(c.lhs, c.rhs) != c.expects {
			t.Errorf("test failed: lhs: %s, rhs: %s, expects: %t", c.lhs, c.rhs, c.expects)
		}
	}
}