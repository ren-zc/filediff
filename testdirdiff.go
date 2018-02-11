package main

import (
	"fmt"
	"github.com/jacenr/filediff/diff"
)

func main() {
	s1 := []string{"a", "b"}
	s2 := []string{}
	s3, s4 := diff.DiffOnly(s1, s2)
	fmt.Println(s3)
	fmt.Println(s4)
}
