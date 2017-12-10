package diff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
)

var srcFile string
var dstFile string
var srcBytes [][]byte
var srcBLen int
var dstBytes [][]byte
var dstBLen int

// src X, dst Y
type point struct {
	x         int
	y         int
	children  []*point
	ready     chan struct{}
	bestChild *point
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
	ready := make(chan struct{})
	p.x = x
	p.y = y
	p.children = children
	p.ready = ready
	p.bestChild = nil
	p.distance = -1
	return p
}

func readFile(file string) ([][]byte, error) {
	fileContent, RErr := ioutil.ReadFile(file)
	if RErr != nil {
		return nil, RErr
	}
	var fileBytes [][]byte
	if runtime.GOOS == "windows" {
		fileBytes = bytes.Split(fileContent, []byte("\r\n")) // windows
	} else {
		fileBytes = bytes.Split(fileContent, []byte{'\n'}) // linux or others
	}
	return fileBytes, nil
}

func initData() error {
	var readErr error
	srcBytes, readErr = readFile(srcFile)
	if readErr != nil {
		return readErr
	}
	srcBLen = len(srcBytes)
	dstBytes, readErr = readFile(dstFile)
	if readErr != nil {
		return readErr
	}
	dstBLen = len(dstBytes)
	return nil
}

func minDistance(pts []*point) *point {
	var min *point
	if _, ok := <-pts[0].ready; !ok {
		min = pts[0]
	}
	for _, v := range pts[1:] {
		if _, ok := <-v.ready; !ok {
			if v.distance < min.distance {
				min = v
			}
		}
	}
	return min
}

func (p *point) getBestPath() {
	if p.x == srcBLen && p.y == dstBLen {
		p.distance = 0
		close(p.ready)
		return
	}
	p.bestChild = minDistance(p.children)
	p.distance = p.bestChild.distance + 1
	close(p.ready)
}

func initGraph() [][]*point {
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

func Diff(dst string, src string) ([][]byte, error) {
	dstFile = dst
	srcFile = src
	initErr := initData()
	if initErr != nil {
		return nil, initErr
	}
	result := make([][]byte, 0, (srcBLen + dstBLen + 1))
	// result = append(result, []byte("@@@ S: src, D: dst. @@@")) // Output head.
	graph := initGraph()
	for _, pl := range graph {
		for _, p := range pl {
			go p.getBestPath()
		}
	}
	pList := make([]*point, 0, (srcBLen + dstBLen + 1))
	var travelPoint func(p *point)
	travelPoint = func(p *point) {
		pList = append(pList, p)
		if p.bestChild == nil {
			return
		} else {
			travelPoint(p.bestChild)
		}
	}
	if _, ok := <-graph[0][0].ready; !ok {
		travelPoint(graph[0][0])
	}
	growResult := func(sr string, i int, byteList [][]byte) {
		s := fmt.Sprintf("%s %d ", sr, i+1)
		b := []byte(s)
		for _, bt := range byteList[i] {
			b = append(b, bt)
		}
		result = append(result, b)
	}
	tmp := pList[0]
	for _, p := range pList[1:] {
		dx := p.x - tmp.x
		dy := p.y - tmp.y
		if dy == 0 {
			growResult("- ", tmp.x, srcBytes)
			tmp = p
			continue
		}
		if dx == 0 {
			growResult("+ ", tmp.y, dstBytes)
			tmp = p
			continue
		}
		growResult("  ", tmp.y, dstBytes)
		tmp = p
	}
	return result, nil
}
