package cpic

import (
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	p := NewParser(`tree:
	-> black
		-> red
		-> red`)
	n := p.Parse()
	n.walk(func(n *node) {
		log.Println(n.ele.(Token), "depth", n.depth)
	})
}
