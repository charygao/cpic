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
		[]string{"expect '>' after '-'"},
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
		lex.name = ltp.name
		i := 0
		for t := lex.token(); t.typ != tEOF; t = lex.token() {
			if ltp.output[i] != t.typ {
				test.Fatalf("%s(wanted) not equal %s in %s\n", tokens[ltp.output[i]], tokens[t.typ], lex.name)
			}
			i++
		}
		for i := 0; i < len(ltp.errs); i++ {
			if ltp.errs[i] != lex.errors[i] {
				test.Fatalf("%s(wanted error) not equal %s in %s\n", ltp.errs[i], lex.errors[i], lex.name)
			}

		}
	}
}

type posTestSuit struct {
	input     string
	positions []position
}

var posTestSuits = []posTestSuit{
	{`tree:
	->black
		->->red
		->red
`,
		[]position{{0, 1}, {0, 5}, {1, 1}, {1, 2}, {1, 4}}},
}

func TestPos(t *testing.T) {

}
