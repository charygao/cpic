package cpic

import (
	//"log"
	"strings"
	"testing"
)

type treeTestSuit struct {
	name   string
	input  string
	output string
	err    string
}

var treeTestSuits = []treeTestSuit{
	{
		"Normal Test 1",
		`tree:
	->black
		->red
		->red`,
		`TREE***
black**
/**\***
red*red
`, ``}, {
		"Normal Test 2",
		`tree:
	-> black
		->red
			->red
				->red
				->red
					->red
					->red
		->red
			->red
				->red
		->red
			->red`, `TREE***************
black**************
/**********\***\***
red*********red*red
/***********/***/**
red*********red*red
/**\********/******
red*red*****red****
****/**\***********
****red*red********
`, ``},
	{"Exception Test 3",
		`tree:
	->red
		->red
			->red
	->black`, ``,
		`unexpected '	' after tree parsed at line 5,column 1`, //do not care output
	},
}

func TestDraw(t *testing.T) {
	defer func() {
		placeHolder = ' '
	}()
	placeHolder = '*'
	for _, suit := range treeTestSuits {
		p := newParser(suit.input)
		n := p.parse()
		if n == nil {
			if suit.err != strings.Join(p.errors, "\n") {
				t.Fatalf("\n%s(wanted error) not equal\n%s", suit.err, strings.Join(p.errors, "\n"))
			}
		} else {
			m := newMatrix(&tree{root: n})
			m.draw()
			if suit.output != m.output() && len(suit.output) > 0 {
				for i := 0; i < len(suit.output); i++ {
					if suit.output[i] != m.output()[i] {
						t.Logf("diff index [%d]: %c %c in %s", i, suit.output[i], m.output()[i], suit.name)
						break
					}
				}
				t.Fatalf("\n%s(wanted output) not equal\n%s in %s", suit.output, m.output(), suit.name)
			}
		}
	}
}
