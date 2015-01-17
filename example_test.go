package cpic

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ExampleTree() {
	sr := strings.NewReader(`tree:
	->black
		->red
		->red
			->black`)
	bs, _ := ioutil.ReadAll(sr)
	out, err := Gen(string(bs))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(out)
	//TREE
	//black
	//|  \
	//red red
	//    |
	//    black
}
func ExampleGraph() {
	sr := strings.NewReader(`graph:
a -> 1 b,2 c,3 d
b -> a,2 c`,
	)
	bs, _ := ioutil.ReadAll(sr)
	out, err := Gen(string(bs))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
	//a-3.000---------------+
	// -2.000--------+      |
	// -1.000+       |      |
	//^      |       |      |
	//|      v       |      |
	//+------b-2.000+|      |
	//              ||      |
	//              vv      |
	//              c       |
	//                      |
	//                      v
	//                      d

}
