package ios

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	str, err := GetVersion()
	if err != nil {
		t.Errorf("error: %s", str)
	}
}