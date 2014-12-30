package cpic

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ExampleGen() {
	sr := strings.NewReader(`tree:
	->black
		->red
		->red
			->black`)
	bs, e := ioutil.ReadAll(sr)
	if e != nil {
		panic(e)
	}
	out, err := Gen(string(bs))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(out)
}
