package bench

import (
	"github.com/jacenr/filediff/diff"
	"testing"
)

func BenchmarkDiff(b *testing.B) {
	for i := 0; i < b.N; i++ {
		diff.Diff("testFile/SrcFile3", "testFile/DstFile3")
	}
}
