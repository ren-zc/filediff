package dirdiffTest

import (
	"github.com/jacenr/filediff/diff"
	"strings"
	"testing"
)

type testDirDiff struct {
	src []string
	dst []string
	rm  string
	add string
}

func TestDiffOnly(t *testing.T) {
	testList := []testDirDiff{
		{[]string{"A", "B", "C", "F", "D", "E"}, []string{"A", "M", "C", "N", "D", "O"}, "BFE", "MNO"},
		// {[]string{"A", "B", "C", "F", "D", "E"}, []string{"A", "B", "C", "F", "D"}, "E", ""},
	}
	for _, v := range testList {
		testRm, testAdd := diff.DiffOnly(v.src, v.dst)
		if strings.Join(testRm, "") != v.rm || strings.Join(testAdd, "") != v.add {
			t.Error("False")
		}
	}
}
