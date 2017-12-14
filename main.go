// This is just for a sample.

package main

import (
	"fmt"
	"github.com/jacenr/filediff/diff"
)

func main() {
	result, _ := diff.Diff("testFile/SrcFile", "testFile/DstFile")
	for _, v := range result {
		fmt.Println(v)
	}
}
