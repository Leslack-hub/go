package main

import (
	"fmt"
	"os"
)

func main() {
	file, e := os.Open("src/meza/test.in")
	if e != nil {
		panic(e)
	}
	var raw, col int
	fmt.Fscan(file, &raw, &col)
	meza := make([][]int, raw)

	for i := range meza {
		meza[i] = make([]int, col)
		for j := range meza[i] {
			fmt.Fscan(file,"%d", &meza[i][j])
		}
	}

	for i := range meza {
		for j := range meza[i] {
			fmt.Printf("%d ", meza[i][j])
		}
		fmt.Println()
	}
}
