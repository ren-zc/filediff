package main

import (
	"fmt"
	"github.com/jacenr/filediff/diff"
)

func main() {
	// diff.Diff()
	var srcFile = "SrcFile"
	var dstFile = "DstFile"
	// a := diff.ReadFile(srcFile)
	// fmt.Println(a)
	graph := diff.InitGraph(srcFile, dstFile)
	// fmt.Println()
	for _, val := range graph {
		for _, val2 := range val {
			fmt.Println(val2)
		}
	}
}
