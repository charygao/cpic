package cpic

import (
	"testing"
)

func TestDraw(t *testing.T) {
	p := NewParser(`tree:
	-> black
		-> red
		-> red`)
	n := p.Parse()
	m := NewMatrix(n)
	m.paintA([][]uint8{{'a', 'b'}, {'d', 'e'}}, 0, 0)
	m.draw()
}
