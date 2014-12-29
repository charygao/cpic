package cpic

import (
	"fmt"
	"strings"
	"testing"
)

type parserTestSuit struct {
	name   string //for debug
	errs   []string
	input  string
	output string
}

var parserTestSuits = []parserTestSuit{
	{"Normal Test 1",
		[]string{},
		`tree:
	->black
		->red
			->red
				->red
			->red
		->red
		->red`,
		`TREE
*black
**red
***red
****red
***red
**red
**red
`},
	{
		"Exception Test 1",
		[]string{"want 2 indents for child node but get 3 at line 3"},
		`tree:
	->black
			->red`,
		``, //do not care output
	},
	{"Exception Test 2",
		[]string{"want a [IDENT], but only get '->' at line 3,column 5"},
		`tree:
	->black
		->->red
		->red
`,
		``, //do not care output
	},
	{"Exception Test 3",
		[]string{"unexpected '	' after tree parsed at line 5,column 1"},
		`tree:
	->red
		->red
			->red
	->black
`, ``, //do not care output
	},
}

func TestParse(t *testing.T) {
	defer func() {
		display = func(v ...interface{}) {}
	}()
	display = func(v ...interface{}) {
		//fmt.Println(v...)
	}
	for _, testSuit := range parserTestSuits {
		p := newParser(testSuit.input)
		n := p.parse()
		var out string
		n.walk(func(n *node) {
			if n.ele != nil {
				out += fmt.Sprintf("%s%s\n", strings.Repeat("*", n.depth), n.ele.(token).lit)
			} else {
				out += fmt.Sprintln("TREE")
			}
		})
		if testSuit.output != out && len(testSuit.output) > 0 {
			t.Fatalf("\n%s(wanted output) not equal\n%s in %s", testSuit.output, out, testSuit.name)
		}
		if n == nil {
			for i := 0; i < len(testSuit.errs); i++ {
				if testSuit.errs[i] != p.errors[i] {
					t.Fatalf("\n%s(wanted error) not equal\n%s in %s", testSuit.errs[i], p.errors[i], testSuit.name)
				}
			}
		}
	}
}
