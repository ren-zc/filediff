package main

import (
	"github.com/jacenr/filediff/diffV2"
)

func main() {
	diffV2.Diff("testFile/SrcFile", "testFile/DstFile")
}
