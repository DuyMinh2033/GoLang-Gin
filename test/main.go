package main

import (
	"fmt"
	"sort"
)

func main() {
	intList := []int{1, 3, 5, 9, 4, 2, 0}
	sort.Sort(sort.Reverse(sort.IntSlice(intList)))//asd
	fmt.Println(intList)
}
