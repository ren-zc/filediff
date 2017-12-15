// This is just for a sample.

package main

import (
	"fmt"
	"github.com/jacenr/filediff/diff"
)

func main() {
	result, _ := diff.Diff("testFile/SrcFile1", "testFile/DstFile1")
	for _, v := range result {
		fmt.Println(v)
	}
}
