//lex is insipred by (Rob Pike's talk)[https://www.youtube.com/watch?v=HxaD_trXwRE]
package cpic

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

//position is the token position in the source text,
//used for error tracing,line and col counts from 0.
type position struct {
	line int //line number
	col  int //colummn number
}

//position 表示的是在源文本中的位置,会把ignore掉的token的位置也考虑进去
//方便错误显示,line 和 col 都是从0开始计数的.

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
	cur token //current scanned token

	src string //source

	pos   int //读取的位置
	start int //起始的位置
	width int //token的宽度

	lineNum int //line counter,count from 0
	colNum  int //columnn counter,count from 1

	errors    []string   //errors stack
	state     stateFn    //状态函数
	tokenChan chan token //token channel
}

const (
	key_words_begin = iota
	//keyword
	kTREE //tree
	kGRAPH
	key_words_end

	token_begin
	//token
	tIDENT //identifier (_|letter)(_|letter|digit)*
	tNUM   //number -?digit*.digit*[E|e]-?digit*

	tRIGHT_ARROW //->
	tLEFT_ARROW  //<-

	tTAG   //\t
	tCOLON //:
	tCOMMA //,

	tOR //|

	tERR //err
	tEOF //eof
	token_end
)

const eof = -1

var tokens = map[int]string{
	kTREE:        "tree",
	kGRAPH:       "graph",
	tRIGHT_ARROW: "[->]",
	tLEFT_ARROW:  "[<-]",
	tTAG:         "[TAG]",
	tIDENT:       "[IDENT]",
	tNUM:         "[NUM]",
	tCOLON:       "[COLON]",
	tCOMMA:       "[COMMA]",
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

//push error message into tracing stack
func (l *lexer) err(e string) {
	l.errors = append(l.errors, e)
}

//format error
func (l *lexer) errf(f string, v ...interface{}) {
	l.errors = append(l.errors, fmt.Sprintf(f, v...))
}

//backup a token
func (l *lexer) backup() {
	l.pos -= l.width
	l.colNum -= l.width
}

//peek a token
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

//ignore a token
func (l *lexer) ignore() {
	//l.colNum 表示源文本对应的位置ignore的时候不会忽略

	l.start = l.pos
}

//emit token
func (l *lexer) emit(typ int) {
	var t = token{
		lit: l.src[l.start:l.pos],
		typ: typ,
		pos: position{l.lineNum, l.colNum - (l.pos - l.start)},
	}
	//fmt.Println("token", t)
	l.tokenChan <- t
	if t.typ == tEOF {
		close(l.tokenChan)
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
	//premature lexical scanning
	//不emit接收方就一直在等待
	l.emit(tEOF)
	return nil
}

func lexUnkown(l *lexer) stateFn {
	//premature lexical scanning
	l.emit(tEOF)
	return nil
}

func lexNum(l *lexer) stateFn {
	r := l.next()
	getNum := false
	if r == '-' { //optional '-'
		r = l.next()
	}
	if r != '.' { //optional digit*
		for unicode.IsDigit(r) {
			getNum = true
			r = l.next()
		}
	}
	if r == '.' { //optional .digit*
		for unicode.IsDigit(r) {
			getNum = true
			r = l.next()
		}
	}
	if r == 'e' || r == 'E' { //optional [E|e]-?digit*
		r = l.next()
		if r == '-' {
			r = l.next()
		}
		for unicode.IsDigit(r) {
			getNum = true
			r = l.next()
		}
	}
	if !getNum {

	}
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
	for ; unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'; r = l.next() {
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
	case unicode.IsLetter(r) || r == '_':
		return lexId
	case unicode.IsDigit(r) || r == '.' || r == '-':
		r2 := l.peek()
		if r == '-' &&
			r2 != '.' && !unicode.IsDigit(r2) && r2 != 'E' && r2 != 'e' {
			goto FL
		}
		l.backup()
		return lexNum
	FL:
		fallthrough
	case r == '-':
		r = l.next()
		if r == '>' {
			l.emit(tRIGHT_ARROW)
		} else {
			l.errf("unexpected '%c' after '-' at line %d,column %d", r, l.lineNum+1, l.colNum)
			return lexError
		}
	case r == '<':
		r = l.next()
		if r == '-' {
			l.emit(tLEFT_ARROW)
		} else {
			l.errf("unexpected '%c' after '<' at line %d,column %d", r, l.lineNum+1, l.colNum)
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
	case r == '|':
		l.emit(tOR)
	case r == '\t':
		l.emit(tTAG)
	case r == ',':
		l.emit(tCOMMA)
	case r == eof:
		return lexEOF
	default:
		l.errf("unkown char '%c' at line %d,column %d", r, l.lineNum+1, l.colNum)
		return lexUnkown
	}
	return lexBegin
}
