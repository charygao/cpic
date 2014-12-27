package cpic

import (
	"container/list"
	"fmt"
)

type Parser struct {
	lexer  *Lexer
	buf    *list.List
	errors []string
}

func NewParser(src string) *Parser {
	return &Parser{
		lexer: NewLexer(src),
		buf:   list.New(),
	}
}

func (p *Parser) Err(e string) {
	p.errors = append(p.errors, e)
}
func (p *Parser) Parse() *node {
	n := &node{}
	token := p.token()
	if token.typ == TREE {
		p.Tree(n)
	}
	return n
}

func (p *Parser) token() Token {
	var tk Token
	if p.buf.Front() == nil {
		tk = p.lexer.Token()
	} else {
		frt := p.buf.Front()
		tk = frt.Value.(Token)
		p.buf.Remove(frt)
	}
	return tk
}
func (p *Parser) backup(t *Token) {
	p.buf.PushFront(t)
}

func (p *Parser) expect(typ int) *Token {
	t := p.token()
	if t.typ != typ {
		p.Err(fmt.Sprintf("line %d,column %d,expect \\t,found %s", t.pos.line+1, t.pos.col+1, t.lit))
	}
	return &t
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
			if t := p.expect(pair.typ); t.typ == pair.typ {
				if pair.f != nil {
					pair.f(t)
				}
			} else {
				break
			}
		}
		return false
	}

	var depth = 0
	if t := p.expect(COLON); t.typ != COLON {
		return n
	}

	parse([]pair{
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

	var tree = func(root *node, depth int) {
		var pairs = make([]pair, depth)
		for i := 0; i < depth; i++ {
			pairs = append(pairs, pair{TAG, func(ts ...Token) { depth++ }})
		}
		newNode := new(node)
		parse(append(pairs, []pair{
			{RIGHT_ARROW, nil},
			{IDENT, func(ts ...Token) {
				newNode.ele = ts[0]
				newNode.depth = depth
				n.add(newNode, LEFT)
			},
			}}...))
		//TODO: parse child tree
		//tagCount = 0

		if t := p.expect(TAG); t.typ == TAG {
			p.backup(t)
			tree(newNode, depth)
		} else {
			p.backup(t)
		}
	}
	tree(n, depth)
	return n
}
