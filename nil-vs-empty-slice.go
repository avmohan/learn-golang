package main

import (
	"fmt"
)

func printSlice(x []int) {
	if x == nil {
		fmt.Printf("NIL!! ")
	}
	fmt.Println(x, len(x), cap(x))
}

func main() {
	a := make([]int, 0)
	var b []int = nil
	printSlice(a)
	printSlice(b)
}
