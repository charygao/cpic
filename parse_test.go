package cpic

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

type parserTestPair struct {
	name   string //for debug
	errs   []string
	input  string
	output string
}

func display(v ...interface{}) {
	fmt.Println(v...)
}

var parserTestPairs = []parserTestPair{
	{"Normal Test",
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
`}, {
		"Exception Test 1",
		[]string{},
		`tree:
	->black
			->red`,
		`TREE
*black
***red`,
	},
}

func TestParse(t *testing.T) {
	for _, testPair := range parserTestPairs {
		p := newParser(testPair.input)
		n := p.parse()
		var out string
		n.walk(func(n *node) {
			if n.ele != nil {
				out += fmt.Sprintf("%s%s\n", strings.Repeat("*", n.depth), n.ele.(token).lit)
			} else {
				out += fmt.Sprintln("TREE")
			}
		})
		if testPair.output != out {
			log.Fatalf("\n%s(wanted output) not equal\n%s", testPair.output, out)
		}
	}
}
