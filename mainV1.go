// diff package usage example.
package main

import (
	"fmt"
	"github.com/jacenr/filediff/diffV1"
)

func main() {
	var srcFile = "testFile/SrcFile2"
	var dstFile = "testFile/DstFile2"
	result, diffErr := diffV1.Diff(srcFile, dstFile)
	if diffErr != nil {
		fmt.Println(diffErr)
	}
	for _, byList := range result {
		fmt.Printf("%s\n", byList)
	}
}
