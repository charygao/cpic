//lex is insipred by (Rob Pike's talk)[https://www.youtube.com/watch?v=HxaD_trXwRE]
package cpic

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

//position is the token position in the source text,
//used for error tracing.
type position struct {
	line int //line number
	col  int //colummn number
}

//position 表示的是在源文本中的位置,会把ignore掉的token的位置也考虑进去
//方便错误显示

//lexical token.
type token struct {
	lit string   //literal value
	typ int      //token type
	pos position //postion
}

//for debug
func (t token) String() string {
	return fmt.Sprintf("<lit:\"%s\",typ:%s,pos line:%d,pos col:%d>", t.lit, tokens[t.typ], t.pos.line, t.pos.col)
}

type stateFn func(l *lexer) stateFn

//lexical parser.
type lexer struct {
	cur token

	src string //source

	pos   int //读取的位置
	start int //起始的位置
	width int //token的宽度

	lineNum int //line counter
	colNum  int //columnn counter

	errors    []string   //errors stack
	state     stateFn    //状态函数
	tokenChan chan token //token channel
}

const (
	key_words_begin = iota
	//keyword
	kTREE
	key_words_end

	token_begin
	//token
	tIDENT
	tRIGHT_ARROW
	tTAG
	tCOLON
	tEOF
	token_end
)

const eof = -1

var tokens = map[int]string{
	kTREE:        "tree",
	tRIGHT_ARROW: "[->]",
	tTAG:         "[TAG]",
	tIDENT:       "[IDENT]",
	tCOLON:       "[COLON]",
	tEOF:         "[EndOfFile]",
}

var keywords map[string]int

func init() {
	keywords = make(map[string]int)
	for k := key_words_begin + 1; k < key_words_end; k++ {
		keywords[tokens[k]] = k
	}
}

//newLexer initiates token channel and go run lexer.run and return lexer.
func newLexer(src string) *lexer {
	l := &lexer{src: src,
		tokenChan: make(chan token),
	}
	go l.run()
	return l
}

//scan next token
func (l *lexer) next() rune {
	if l.pos >= len(l.src) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.src[l.pos:])
	l.width = w
	l.pos += l.width
	l.colNum += l.width
	return r
}

func (l *lexer) consume() {
}

//backup a token
func (l *lexer) backup() {
	l.pos -= l.width
	l.colNum -= l.width
}

//ignore a token
func (l *lexer) ignore() {
	//l.colNum 表示源文本对应的位置ignore的时候不会忽略

	l.start = l.pos
}

//emit token
func (l *lexer) emit(typ int) {
	l.tokenChan <- token{
		lit: l.src[l.start:l.pos],
		typ: typ,
		pos: position{l.lineNum, l.colNum - (l.pos - l.start)},
	}
	l.start = l.pos

}

//read token
func (l *lexer) token() token {
	token := <-l.tokenChan
	l.cur = token
	return token
}

//main consumer routine
func (l *lexer) run() {
	for l.state = lexBegin; l.state != nil; {
		l.state = l.state(l)
	}
}

//------------------------------------sate function----------------------------------
//error handling
//TODO:sync errors.
func lexError(l *lexer) stateFn {
	return lexBegin
}

//end of scanning
func lexEOF(l *lexer) stateFn {
	l.emit(tEOF)
	return nil
}

//parse id
func lexId(l *lexer) stateFn {
	r := l.next()
	for ; unicode.IsLetter(r); r = l.next() {
	}
	l.backup()
	v := l.src[l.start:l.pos]
	if tok, ok := keywords[v]; ok {
		l.emit(tok)
		return lexBegin
	}
	l.emit(tIDENT)
	return lexBegin
}

//main parse entry
func lexBegin(l *lexer) stateFn {
	switch r := l.next(); {
	case unicode.IsLetter(r):
		return lexId
	case r == '-':
		r = l.next()
		if r == '>' {
			l.emit(tRIGHT_ARROW)
		} else {
			l.backup()
			l.errors = append(l.errors, "expect '>' after '-'")
			return lexError
		}
	case r == ' ':
		l.ignore()
	case r == '\n':
		l.ignore()
		l.lineNum++
		l.colNum = 0
	case r == ':':
		l.emit(tCOLON)
	case r == '\t':
		l.emit(tTAG)
	case r == eof:
		return lexEOF
	}
	return lexBegin
}
