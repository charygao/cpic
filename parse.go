//TODO:add more test cases.
package cpic

import (
	"container/list"
	"fmt"
	"strings"
	//"log"
)

//parser is the the syntax parser
type parser struct {
	lexer  *lexer
	buf    *list.List
	errors []string
}

//Newparser initiates parser to parse src.
func newParser(src string) *parser {
	return &parser{
		lexer: newLexer(src),
		buf:   list.New(),
	}
}

//place holder function
func __placeholder() {
	//log.Println("log")
}

//Err appends error to error stack.
func (p *parser) err(e string) {
	p.errors = append(p.errors, e)
}

//format error
func (p *parser) errf(f string, v ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(f, v...))
}

//parse is the main parsing entry.
func (p *parser) parse() *node {
	n := &node{}
	token := p.token()
	if token.typ == kTREE {
		//parse a tree
		n = p.Tree(n)
	}
	return n
}

//token gets token from buf,if not,waits for the chanell to recevie new one.
func (p *parser) token() token {
	var tk token
	if p.buf.Front() == nil {
		if len(p.lexer.errors) > 0 {
			//come across errors,stop parsing.
			panic(strings.Join(p.lexer.errors, "\n"))
		}
		tk = p.lexer.token()
	} else {
		frt := p.buf.Front()
		tk = frt.Value.(token)
		p.buf.Remove(frt)
	}
	return tk
}

//peek a token
func (p *parser) peek() token {
	var t = p.token()
	p.backup(t)
	return t
}

func (p *parser) backup(t token) {
	p.buf.PushFront(t)
}

//side effect:error handling
func (p *parser) expect(typ int) bool {
	t := p.token()
	if t.typ != typ {
		p.err(fmt.Sprintf("line %d,column %d,expect \\t,found %s", t.pos.line+1, t.pos.col+1, t.lit))
		return false
	}
	return true
}

//parse a run of tokens,if not backup all the parsed the tokens in this routine
func (p *parser) run(typs []int) bool {
	var tkbuf []token
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
func (p *parser) foresee(typs []int) bool {
	var ok = true
	var tkbuf []token
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

//parse tree.
func (p *parser) Tree(n *node) *node {
	//keyword tree has been parsed

	//token action pair
	type pair struct {
		typ int
		f   func(token ...token)
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
	if !p.expect(tCOLON) {
		return n
	}

	var tree func(root *node) *node
	tree = func(root *node) *node {
		//log.Println("depth", root.depth)
		//希望看到长度为depth的[TAG]token
		var tags []int
		for i := 0; i < root.depth+1; i++ {
			tags = append(tags, tTAG)
		}

		//每次看到长度为depth的[TAG]token解析一个结点
		//want "depth" tags
		for p.run(tags) {
			//log.Println("run tags ", root.depth)
			//解析一个结点
			//parse a new node
			newNode := new(node)
			ok := parse(
				[]pair{
					{tRIGHT_ARROW, nil},
					{tIDENT, func(ts ...token) {
						newNode.ele = ts[0]
						root.add(newNode, rRIGHT)

						//log.Println("add node", newNode.depth, ts[0])
						if root.ele != nil {
							//log.Println("parent node", root.ele.(token))
						}
					}},
				})
			if ok {
				////log.Println("parse node ok")
				//如果下一层的tags数比本层多一个,就解析子树
				//if next tags - current tags = 1,parse child tree
				if p.foresee(append(tags, tTAG)) {
					//log.Println("parse children")
					tree(newNode)
				}
				//parse siblings
				if p.foresee(tags) {
					continue
				}
				//log.Println("parse siblings")
				//条件不成立的话解析掉的tokens都会重新退回到缓存里面
			} else {
				p.errf("want -> identifer,but get %s at line %d,column %d", p.peek().lit, p.peek().pos.line, p.peek().pos.col)
				return nil
				//log.Println("parse node not ok")
			}
		}
		return n
	}
	return tree(n)
}
