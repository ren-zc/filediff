package diff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var srcFile = "SrcFile"
var dstFile = "DstFile"
var lg *log.Logger
var wg sync.WaitGroup

func init() {
	lg = log.New(os.Stdout, "diff ", log.Lshortfile)
}

// src X, dst Y
type point struct {
	x         int
	y         int
	children  []*point
	ready     chan struct{}
	bestChild int
	distance  int
}

func (p *point) String() string {
	var s string
	for _, v := range p.children {
		s += fmt.Sprintf(" %d,%d", v.x, v.y)
	}
	return fmt.Sprintf("%d,%d", p.x, p.y) + "\t" + s
}

func initPoint(x int, y int) *point {
	p := new(point)
	children := make([]*point, 0, 3)
	bestChild := int(-1)
	ready := make(chan struct{})
	p.x = x
	p.y = y
	p.children = children
	p.ready = ready
	p.bestChild = bestChild
	p.distance = int(-1)
	return p
}

func readFile(file string) [][]byte {
	fileContent, RErr := ioutil.ReadFile(file)
	if RErr != nil {
		lg.Fatalln(RErr)
	}
	fileBytes := bytes.Split(fileContent, []byte("\r\n")) // windows
	// fileBytes := bytes.Split(fileContent, []byte{'\n'})   // linux
	return fileBytes
}

func InitGraph(src string, dst string) [][]*point {
	srcBytes := readFile(src)
	srcBLen := len(srcBytes)
	dstBytes := readFile(dst)
	dstBLen := len(dstBytes)
	theSame := make(map[int][]int)
	graph := make([][]*point, 0, (srcBLen + 1))
	for i := 0; i <= srcBLen; i++ {
		graphY := make([]*point, 0, (dstBLen + 1))
		for j := 0; j <= dstBLen; j++ {
			p := initPoint(i, j)
			graphY = append(graphY, p)
		}
		graph = append(graph, graphY)
	}
	for i, srcB := range srcBytes {
		dstBList := []int{}
		for j, dstB := range dstBytes {
			if bytes.Equal(srcB, dstB) {
				dstBList = append(dstBList, j)
			}
		}
		if len(dstBList) != 0 {
			theSame[i] = dstBList
		}
	}
	for i, srcX := range graph {
		for j, _ := range srcX {
			if i < srcBLen {
				graph[i][j].children = append(graph[i][j].children, graph[i+1][j])
			}
			if j < dstBLen {
				graph[i][j].children = append(graph[i][j].children, graph[i][j+1])
			}
			if v, ok := theSame[i]; ok {
				for _, vv := range v {
					if j == vv {
						graph[i][j].children = append(graph[i][j].children, graph[i+1][j+1])
						break
					}
				}
			}
		}
	}
	return graph
}

func Diff() {
	fmt.Println("Test.")
}
