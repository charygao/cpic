//This is a command line interface package for cpic.
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ggaaooppeenngg/cpic"
)

func main() {
	sr := bufio.NewReader(os.Stdin)
	bs, e := ioutil.ReadAll(sr)
	if e != nil {
		panic(e)
	}
	fmt.Printf("%s\n", bs)
	out := cpic.Gen(string(bs))
	fmt.Println(out)
}
