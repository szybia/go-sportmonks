package sportmonks

import (
	"testing"
)

func TestGet(t *testing.T) {
	_, err := Get("testing", "", 0, false)
	if err == nil {
		t.Error("Get function does not check for empty API Token")
	}
	SetAPIToken("1234")

	_, err = Get("", "", 0, false)
	if err == nil {
		t.Errorf("Get function accepts empty endpoint")
	}
}
