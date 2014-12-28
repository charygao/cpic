package cpic

import (
	"container/list"
	"fmt"
	"log"
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

//place holder function
func __placeholder() {
	log.Println("log")
}

func (p *Parser) Err(e string) {
	p.errors = append(p.errors, e)
}
func (p *Parser) Parse() *node {
	n := &node{}
	token := p.token()
	if token.typ == TREE {
		//parse a tree
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

func (p *Parser) backup(t Token) {
	p.buf.PushFront(t)
}

//side effect:error handling
func (p *Parser) expect(typ int) bool {
	t := p.token()
	if t.typ != typ {
		p.Err(fmt.Sprintf("line %d,column %d,expect \\t,found %s", t.pos.line+1, t.pos.col+1, t.lit))
		return false
	}
	return true
}

//parse a run of tokens,if not backup all the parsed the tokens in this routine
func (p *Parser) run(typs []int) bool {
	var tkbuf []Token
	for _, typ := range typs {
		t := p.token()
		tkbuf = append(tkbuf, t)
		if typ != t.typ {
			for i := len(tkbuf); i > 0; i-- {
				p.backup(tkbuf[i-1])
			}
			return false
		}
	}
	return true
}

//foresee run of tokens
func (p *Parser) foresee(typs []int) bool {
	var ok = true
	var tkbuf []Token
	for _, typ := range typs {
		t := p.token()
		tkbuf = append(tkbuf, t)
		if typ != t.typ {
			ok = false
			break
		}
	}
	for i := len(tkbuf); i > 0; i-- {
		p.backup(tkbuf[i-1])
	}
	return ok
}

func (p *Parser) Tree(n *node) *node {
	//keyword tree has been parsed

	//token action pair
	type pair struct {
		typ int
		f   func(token ...Token)
	}

	//if parsing unwanted token,break the parsing and return false
	var parse = func(pairs []pair) bool {
		for _, pair := range pairs {
			if t := p.token(); t.typ == pair.typ {
				if pair.f != nil {
					pair.f(t)
				}
			} else {
				return false
			}
		}
		return true
	}

	//want COLON after keyword tree
	if !p.expect(COLON) {
		return n
	}

	var tree func(root *node)
	tree = func(root *node) {
		log.Println("depth", root.depth)
		//希望看到长度为depth的[TAG]token
		var tags []int
		for i := 0; i < root.depth+1; i++ {
			tags = append(tags, TAG)
		}

		//每次看到长度为depth的[TAG]token解析一个结点
		for p.run(tags) {
			log.Println("run tags ", root.depth)
			//解析一个结点
			newNode := new(node)
			ok := parse(
				[]pair{
					{RIGHT_ARROW, nil},
					{IDENT, func(ts ...Token) {
						newNode.ele = ts[0]
						root.add(newNode, LEFT)

						log.Println("add node", newNode.depth, ts[0])
						if root.ele != nil {
							log.Println("parent node", root.ele.(Token))
						}
					}},
				})
			if ok {
				//log.Println("parse node ok")
				//如果下一层的tags数比本层多一个,就解析子树
				if p.foresee(append(tags, TAG)) {
					log.Println("parse children")
					tree(newNode)
				} else {
					log.Println("parse siblings")
				}
				//条件不成立的话解析掉的tokens都会重新退回到缓存里面
			} else {
				log.Println("parse node not ok")
			}
		}
	}
	tree(n)
	return n
}
