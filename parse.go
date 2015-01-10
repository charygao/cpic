//TODO:add more test cases.
package cpic

import (
	"container/list"
	"fmt"
	//"log"
	"strconv"
	"strings"

	"github.com/ggaaooppeenngg/util/container"
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
func (p *parser) parse() painter {
	token := p.token()
	if token.typ == kTREE {
		tree := newTree()
		//parse a tree
		p.Tree(tree)
		return tree
	}
	if token.typ == kGRAPH {
		graph := newGraph()
		//parse graph
		p.Graph(graph)
		return graph
	}
	token = p.token()
	//TODO:varibal tree
	if token.typ != tEOF {
		p.errf("unexpected '%s' after end of tree parsing at line %d,column %d", token.lit, token.pos.line+1, token.pos.col+1)
	}
	//fmt.Println("return nil")
	return nil
}

//token gets token from buf,if not,waits for the chanel to recevie new one.
func (p *parser) token() token {
	var tk token
	if p.buf.Front() == nil {
		if len(p.lexer.errors) > 0 {
			//come across errors,stop parsing.
			//TODO:remove panic to handle error,use nil token insead.
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

//token action pair
type pair struct {
	typ int
	f   func(token ...token)
}

//if parsing unwanted token,break the parsing and return false
func (p *parser) parsePairs(pairs []pair) (t *token, typ int) {
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

//graph 的语法.
//
//图不可以赋值.
//graph:
//a->1 b,10 c,1 d,2 e,100 f//不可以重复,没有的会自动创建,但是不能有平行边.
//b->f,2 c,e,f  //没有出现的,自动添加,没有默认是0,
//only decimal float and int allowed in weight.
func (p *parser) Graph(n *graph) {
	//keyword 'graph' has been parsed.
	//expect a colon
	p.expect(tCOLON)
	//parsing----
	p.parseGraph(n)
}

func (p *parser) parseEdge(vtx *container.Vertex, g *graph) *container.Vertex {
	token := p.token()
	var weight float64 = 0.0
	var numParsed = false
	if token.typ == tNUM {
		numParsed = true
		weight, _ = strconv.ParseFloat(token.lit, 64)
	}
	if numParsed {
		token = p.token()
	}
	if token.typ == tIDENT {
		v := g.FindVertexById(token.lit)
		if v == nil {
			newVtx := &container.Vertex{token.lit, weight}
			g.Connect(vtx, newVtx, weight)
			return newVtx
		} else {
			if g.EdgeExists(vtx, v) {
				p.errf("edge from %d to %d has exists,can not reset it", vtx.Id, v.Id)
			} else {
				g.Connect(vtx, v, weight)
				return v
			}
		}
	} else {
		p.errf("want a identifier but get %d at line %d,column %d", token.lit, token.pos.line+1, token.pos.col+1)
	}
	return nil
}

//graph parsing
func (p *parser) parseGraph(g *graph) *graph {
	if p.foresee([]int{tIDENT, tRIGHT_ARROW}) {
		p.errf("emty graph is no allowed")
	}
	for p.foresee([]int{tIDENT, tRIGHT_ARROW}) {
		var vtx *container.Vertex
		p.parsePairs([]pair{
			{tIDENT, func(ts ...token) {
				tk := ts[0]
				vtx = &container.Vertex{Id: tk.lit}
				g.AddVertex(vtx)
			}},
			{tRIGHT_ARROW, nil},
		})
		if p.parseEdge(vtx, g) == nil {
			return g
		}
		for {
			if !p.expect(tCOMMA) {
				return g
			}
			if p.parseEdge(vtx, g) == nil {
				return g
			}
		}
	}
	return g
}

//parse tree.
//tree 的解析可以更复杂一点比如 赋值
//a := tree:
//		->black
//		->red
//		->red
//
//b := tree:
//		->black
//			->a
//
func (p *parser) Tree(t *tree) {
	//keyword tree has been parsed

	//want COLON after keyword tree
	p.expect(tCOLON)
	t.root = p.tree(t.root)
}

func (p *parser) tree(root *node) *node {
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
		errToken, typ := p.parsePairs(
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
			//log.Println("parse node ok")
			//如果下一层的tags数比本层多一个,就解析子树
			//if next tags - current tags = 1,parse child tree
			if p.foresee(append(tags, tTAG)) {
				//log.Println("parse children")
				p.tree(newNode)
			}
			//can parse siblings ?
			if p.foresee(tags) {
				//第一层的结点唯一
				//first level node only allows one.
				if newNode.depth == 1 {
					break
				}
			} else {
				//如果不能解析下一次兄弟结点这层的解析结束,不作错误处理,
				//在树外继续解析.
				break
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
				p.errf("want a %s, but only get '%s' at line %d,column %d", tokens[typ], errToken.lit, errToken.pos.line+1, errToken.pos.col+1)
				return nil
			}
		}

	}
	//tags not enough or come across other tokens
	return root
}
