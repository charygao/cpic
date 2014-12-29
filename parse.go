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

//for debug
var display = func(v ...interface{}) {
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
	var parse = func(pairs []pair) (t *token, typ int) {
		for _, pair := range pairs {
			if t := p.token(); t.typ == pair.typ {
				display(t)
				if pair.f != nil {
					pair.f(t)
				}
			} else {
				display(t)
				return &t, pair.typ
			}
		}
		return nil, 0
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
		//first parsed child
		var firstChild = true
		//每次看到长度为depth的[TAG]token解析一个结点
		//want "depth" tags
		for p.run(tags) {
			//log.Println("run tags ", root.depth)
			//解析一个结点
			//parse a new node
			newNode := new(node)
			errToken, typ := parse(
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
			//如果本身的结点解析成功,再尝试解析子结点
			if errToken == nil {
				firstChild = false
				////log.Println("parse node ok")
				//如果下一层的tags数比本层多一个,就解析子树
				//if next tags - current tags = 1,parse child tree
				if p.foresee(append(tags, tTAG)) {
					//log.Println("parse children")
					tree(newNode)
				}
				//parse siblings
				//如果不能解析下一次的兄弟结点
				if !p.foresee(tags) {
					p.errf("parse tree error,please check wether indent number match or there are any other characters  at line %d", p.peek().lit, p.peek().pos.line+1)
				}
				//log.Println("parse siblings")
				//条件不成立的话解析掉的tokens都会重新退回到缓存里面
			} else {
				//如果本身是父结点的第一个子节点还可能是缩进过多.
				if firstChild && errToken.typ == tTAG {
					p.errf("want %d indents for child node but get %d at line %d", root.depth+1, root.depth+2, errToken.pos.line+1)
					return nil
				} else {
					//如果本身的结点解析不成功,报本身的错误.
					p.errf("want a %s, but only get '%s' at line %d,column %d", tokens[typ], errToken.lit, errToken.pos.line+1, errToken.pos.col)
					return nil
				}
			}

		}
		//tags not enough or come across other tokens
		return n
	}
	return tree(n)
}
