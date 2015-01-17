package cpic

import (
	//"log"
	"strings"
	"testing"
)

func TestTreeDraw(t *testing.T) {
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
			`black**
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
			->red`, `black**************
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
			`unexpected '	' after end of tree parsing at line 5,column 1`, //do not care output
		},
	}

	defer func() {
		placeHolder = ' '
	}()
	placeHolder = '*'
	for _, suit := range treeTestSuits {
		p := newParser(suit.input)
		painter := p.parse()
		if painter == nil {
			if suit.err != strings.Join(p.errors, "\n") {
				t.Fatalf("\n%s(wanted error) not equal\n%s", suit.err, strings.Join(p.errors, "\n"))
			}
		} else {
			m := newMatrix(painter)
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

func TestDrawGraph(t *testing.T) {
	type testPair struct {
		input string
	}
	var tests = []testPair{
		/*{`graph:
		a->1 b`,
				}, {`graph:
		a->1 b,2 c`,
				}, {`graph:
		a->b,c`,
				}, */{`graph:
a->1 b,2 c
b->1 a,2 c
c->1 d`,
		}}
	defer func() {
		placeHolder = ' '
	}()
	placeHolder = '*'
	for i, test := range tests {
		t.Log("test", i)
		p := newParser(test.input)
		g := p.parse()
		m := newMatrix(g)
		m.draw()
		t.Log("\n" + m.output())
	}
}
