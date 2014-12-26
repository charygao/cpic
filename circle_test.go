package cpic

import (
	"strings"
	"testing"
)

func TestDrawCircle(t *testing.T) {
	c := &Circle{}
	println(strings.Join(c.Content("hello word"), "\n"))
}
