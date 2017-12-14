package main

import (
	"fmt"
	"github.com/jacenr/filediff/diff"
)

func main() {
	result, _ := diff.Diff("testFile/SrcFile5", "testFile/DstFile5")
	for _, v := range result {
		fmt.Println(v)
	}
}
