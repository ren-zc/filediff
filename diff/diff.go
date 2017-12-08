package diff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "sync"
)

var srcFile = "SrcFile"
var dstFile = "DstFile"
var lg *log.Logger

// var wg sync.WaitGroup

var srcBytes [][]byte
var srcBLen int
var dstBytes [][]byte
var dstBLen int

func init() {
	lg = log.New(os.Stdout, "diff ", log.Lshortfile)
	srcBytes = readFile(srcFile)
	srcBLen = len(srcBytes)
	dstBytes = readFile(dstFile)
	dstBLen = len(dstBytes)
}

// src X, dst Y
type point struct {
	x         int
	y         int
	children  []*point
	ready     chan struct{}
	bestChild *point
	distance  int
}

func (p *point) getBestPath() {
	// fmt.Println()
	// fmt.Println(p)
	if p.x == srcBLen && p.y == dstBLen {
		p.distance = 0
		close(p.ready)
		return
		// return p
	}
	childNum := len(p.children)
	if childNum == 0 {
		lg.Fatalf("Point (x:%d,y:%d) childern number is 0, please check", p.x, p.y)
	}
	// fmt.Println(p.children)
	p.bestChild = minDistance(p.children)
	p.distance = p.bestChild.distance + 1
	close(p.ready)
	// return nil
}

func minDistance(pts []*point) *point {
	// fmt.Println(len(pts))
	var min *point
	// fmt.Println("1")
	if _, ok := <-pts[0].ready; !ok {
		// fmt.Println("in")
		min = pts[0]
		// fmt.Println(min)
	}
	// fmt.Println("2")
	for i, v := range pts {
		if i == 0 {
			continue
		}
		if _, ok := <-v.ready; !ok {
			// fmt.Println("ok")
			// fmt.Printf("%d,%v,%v\n", i, v, min)
			// fmt.Printf("%d,%d\n", v.distance, min.distance)
			if v.distance < min.distance {
				min = v
			}
		}
	}
	// fmt.Println(min)
	return min
}

func (p *point) String() string {
	var s string
	// for _, v := range p.children {
	// 	s += fmt.Sprintf(" %d,%d", v.x, v.y)
	// }
	return fmt.Sprintf("%d,%d", p.x, p.y) + "\t" + s
}

func initPoint(x int, y int) *point {
	p := new(point)
	children := make([]*point, 0, 3)
	// bestChild := int(-1)
	ready := make(chan struct{})
	p.x = x
	p.y = y
	p.children = children
	p.ready = ready
	p.bestChild = nil
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

func InitGraph() [][]*point {
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
						// fmt.Println(graph[i+1][j+1])
						break
					}
				}
			}
		}
	}
	return graph
}

func Diff() {
	// fmt.Println("Test.")
	graph := InitGraph()
	for _, vx := range graph {
		for _, vy := range vx {
			// fmt.Println(vy)
			go vy.getBestPath()
		}
	}
	if _, ok := <-graph[0][0].ready; !ok {
		printPoint(graph[0][0])
	}
}

func printPoint(p *point) {
	fmt.Println(p)
	if p.bestChild == nil {
		return
	} else {
		printPoint(p.bestChild)
	}
}
