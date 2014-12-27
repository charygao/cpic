//insipred by Rob Pike lexer
package cpic

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

//lexical parser
//表示的是在源文本中的位置,会把ignore掉的token的位置也考虑进去
//方便错误显示
type Position struct {
	line int //line number
	col  int //colummn number
}

type Token struct {
	lit string   //literal value
	typ int      //token type
	pos Position //postion
}

//for debug
func (t Token) String() string {
	return fmt.Sprintf("<lit:\"%s\",typ:%s,pos line:%d,pos col:%d>", t.lit, tokens[t.typ], t.pos.line, t.pos.col)
}

type stateFn func(l *Lexer) stateFn

type Lexer struct {
	cur Token

	src string //source

	pos   int //读取的位置
	start int //起始的位置
	width int //token的宽度

	lineNum int //line counter
	colNum  int //columnn counter

	errors    []string   //errors stack
	state     stateFn    //状态函数
	tokenChan chan Token //token channel
}

const (
	key_words_begin = iota
	TREE
	key_words_end

	token_begin
	IDENT
	RIGHT_ARROW
	TAG
	COLON
	EOF
	token_end
)

const eof = -1

var tokens = map[int]string{
	TREE:        "tree",
	RIGHT_ARROW: "[->]",
	TAG:         "[TAG]",
	IDENT:       "[IDENT]",
	COLON:       "[COLON]",
	EOF:         "[EndOfFile]",
}
var keywords map[string]int

func init() {
	keywords = make(map[string]int)
	for k := key_words_begin + 1; k < key_words_end; k++ {
		keywords[tokens[k]] = k
	}
}

func NewLexer(src string) *Lexer {
	l := &Lexer{src: src,
		tokenChan: make(chan Token),
	}
	go l.run()
	return l
}

//scan next token
func (l *Lexer) next() rune {
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

func (l *Lexer) consume() {
}

//backup a token
func (l *Lexer) backup() {
	l.pos -= l.width
	l.colNum -= l.width
}

//ignore a token
func (l *Lexer) ignore() {
	//l.colNum 表示源文本对应的位置ignore的时候不会忽略
	l.start = l.pos
}

//emit token
func (l *Lexer) emit(typ int) {
	l.tokenChan <- Token{
		lit: l.src[l.start:l.pos],
		typ: typ,
		pos: Position{l.lineNum, l.colNum - (l.pos - l.start)},
	}
	l.start = l.pos

}

//read token
func (l *Lexer) Token() Token {
	token := <-l.tokenChan
	l.cur = token
	return token
}

//main consumer routine
func (l *Lexer) run() {
	for l.state = lexBegin; l.state != nil; {
		l.state = l.state(l)
	}
}

//------------------------------------sate function----------------------------------
//error handling
func lexError(l *Lexer) stateFn {
	return lexBegin
}

//end of scanning
func lexEOF(l *Lexer) stateFn {
	l.emit(EOF)
	return nil
}

//parse id
func lexId(l *Lexer) stateFn {
	r := l.next()
	for ; unicode.IsLetter(r); r = l.next() {
	}
	l.backup()
	v := l.src[l.start:l.pos]
	if tok, ok := keywords[v]; ok {
		l.emit(tok)
		return lexBegin
	}
	l.emit(IDENT)
	return lexBegin
}

//main parse entry
func lexBegin(l *Lexer) stateFn {
	switch r := l.next(); {
	case unicode.IsLetter(r):
		return lexId
	case r == '-':
		r = l.next()
		if r == '>' {
			l.emit(RIGHT_ARROW)
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
		l.emit(COLON)
	case r == '\t':
		l.emit(TAG)
	case r == eof:
		return lexEOF
	}
	return lexBegin
}
