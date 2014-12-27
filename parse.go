package cpic

import (
	"fmt"
)

type Parser struct {
	lexer  *Lexer
	errors []string
}

func (p Parser) String() string {
	return fmt.Sprintf("hello world")
}
func NewParser(src string) *Parser {
	return &Parser{
		lexer: NewLexer(src),
	}
}

func (p *Parser) Err(e string) {
	p.errors = append(p.errors, e)
}
func (p *Parser) Parse() *node {
	n := &node{}
	token := p.lexer.Token()
	if token.typ == TREE {
		p.Tree(n)
	}
	return n
}
func (p *Parser) Tree(n *node) *node {
	/*
		var sync = func() {

		}
	*/
	type pair struct {
		typ int
		f   func(token ...Token)
	}

	var parse = func(pairs []pair) (hasErr bool) {
		for _, pair := range pairs {
			t := p.lexer.Token()
			if t.typ != pair.typ {
				p.Err(fmt.Sprintf("line %d,column %d,expect \\t,found %s", t.pos.line+1, t.pos.col+1, t.lit))
				return true
			} else {
				if pair.f != nil {
					pair.f(t)
				}
			}
		}
		return false
	}

	var depth = 0
	parse([]pair{
		{COLON,
			nil},
		{TAG,
			func(ts ...Token) {
				depth++
			}},
		{RIGHT_ARROW,
			nil},
		{IDENT,
			//depth 0
			func(ts ...Token) {
				n.ele = ts[0]
			}},
	})
	return n
}
