package cpic

import (
	"testing"
)

type lexTestPair struct {
	name   string //for debug
	errs   []string
	input  string
	output []int
}

var lexTestPairs = []lexTestPair{
	{
		"Normal Test",
		[]string{},
		`tree:
	-> black
		-> red
		-> red`,
		[]int{kTREE, tCOLON, tTAG, tRIGHT_ARROW, tIDENT,
			tTAG, tTAG, tRIGHT_ARROW, tIDENT,
			tTAG, tTAG, tRIGHT_ARROW, tIDENT}},
	{
		"Exception Test 1",
		[]string{"expect '>' after '-' at line 2,column 3"},
		`tree:
a-<`,
		[]int{kTREE, tCOLON, tIDENT},
	},
	{
		"Exception Test 2",
		[]string{"unkown char '*' at line 2,column 1"},
		`tree
*-> black`,
		[]int{kTREE},
	}, {
		"Exception Test 3",
		[]string{"unkown char '\\' at line 5,column 3"},
		`

	->black
		-> red
		\`,
		[]int{tTAG, tRIGHT_ARROW, tIDENT, tTAG, tTAG, tRIGHT_ARROW, tIDENT, tTAG, tTAG},
	},
}

func TestLex(test *testing.T) {
	for _, ltp := range lexTestPairs {
		lex := newLexer(ltp.input)
		i := 0
		for t := lex.token(); t.typ != tEOF; t = lex.token() {
			if ltp.output[i] != t.typ {
				test.Fatalf("%s(wanted) not equal %s in %s\n", tokens[ltp.output[i]], tokens[t.typ], ltp.name)
			}
			i++
		}
		for i := 0; i < len(ltp.errs); i++ {
			if ltp.errs[i] != lex.errors[i] {
				test.Fatalf("%s(wanted error) not equal %s in %s\n", ltp.errs[i], lex.errors[i], ltp.name)
			}

		}
	}
}

type posTestSuit struct {
	name      string
	input     string
	positions []position
}

var posTestSuits = []posTestSuit{
	{
		`Normal Test`,
		`tree:
	->black
		->->red
		->red
`,
		[]position{{0, 0}, {0, 4},
			{1, 0}, {1, 1}, {1, 3},
			{2, 0}, {2, 1}, {2, 2}, {2, 4}, {2, 6},
			{3, 0}, {3, 1}, {3, 2}, {3, 4},
		}},
}

func TestPos(test *testing.T) {
	for _, suit := range posTestSuits {
		l := newLexer(suit.input)
		i := 0
		for t := l.token(); t.typ != tEOF; t = l.token() {
			if t.pos != suit.positions[i] {
				test.Fatalf("want pos [%d] line %d,col %d,but get line %d,col %d(token is %v )", i, suit.positions[i].line, suit.positions[i].col, t.pos.line, t.pos.col, t)
			}
			i++
		}
	}

}
