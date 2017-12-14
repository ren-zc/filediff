package diffV2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
)

var srcFile [][]byte
var dstFile [][]byte
var srcLen int
var dstLen int

// src:x, dst: y
type point struct {
	x int
	y int
}

func (p *point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func newpoint(x int, y int) *point {
	p := new(point)
	p.x = x
	p.y = y
	return p
}

var path map[*point][]*point
var newed map[string]*point

func init() {
	path = map[*point][]*point{}
	newed = map[string]*point{}
}

func checkNew(x, y int) *point {
	xyStr := strconv.Itoa(x) + strconv.Itoa(y)
	v, ok := newed[xyStr]
	if !ok {
		p := newpoint(x, y)
		newed[xyStr] = p
		return p
	} else {
		return v
	}
}

func scanPath(p *point) []*point {
	shortPath := []*point{}
	xlimit := srcLen
	ylimit := dstLen
	for x, y := p.x+1, p.y+1; x < xlimit && y < ylimit; x, y = x+1, y+1 {
		if bytes.Equal(srcFile[x], dstFile[y]) {
			pn := checkNew(x, y)
			shortPath = append(shortPath, pn)
			return shortPath
		}
		for i := x + 1; i < xlimit; i++ {
			if bytes.Equal(srcFile[i], dstFile[y]) {
				xlimit = i
				// pi := newpoint(i, y)
				pi := checkNew(i, y)
				shortPath = append(shortPath, pi)
				break
			}
		}
		for j := y + 1; j < ylimit; j++ {
			if bytes.Equal(dstFile[j], srcFile[x]) {
				ylimit = j
				// pj := newpoint(x, j)
				pj := checkNew(x, j)
				shortPath = append(shortPath, pj)
				break
			}
		}
	}
	return shortPath
}

func getPath(p *point) {
	if _, ok := path[p]; ok {
		return
	}
	ps := scanPath(p)
	if len(ps) == 0 {
		return
	}
	path[p] = ps
	for _, pn := range ps {
		getPath(pn)
	}
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
		fileBytes = bytes.Split(fileContent, []byte{'\n'}) // linux and others
	}
	return fileBytes, nil
}

func getMostDepth(p *point) []*point {
	children, ok := path[p]
	if !ok {
		pl := []*point{}
		pl = append(pl, p)
		return pl
	}
	depth := 0
	var pl []*point
	for _, v := range children {
		plv := getMostDepth(v)
		if length := len(plv); length > depth {
			depth = length
			pl = plv
		}
	}
	pl = append(pl, p)
	return pl
}

func Diff(src string, dst string) {
	srcFile, _ = readFile(src)
	dstFile, _ = readFile(dst)
	srcLen = len(srcFile)
	dstLen = len(dstFile)
	pTmp := newpoint(-1, -1)
	getPath(pTmp)
	// for k, v := range path {
	// 	fmt.Printf("%v\t%v\n", k, v)
	// }
	pathPoint := getMostDepth(pTmp)
	// fmt.Println(pathPoint)

	result := []string{}
	var str string
	pOne := newpoint(0, 0)
	getResult := func(pOne, pPoint *point) {
		for j := pOne.x; j < pPoint.x; j++ {
			str = fmt.Sprintf("- %s", srcFile[j])
			result = append(result, str)
		}
		for j := pOne.y; j < pPoint.y; j++ {
			str = fmt.Sprintf("+ %s", dstFile[j])
			result = append(result, str)
		}
	}
	for i := len(pathPoint) - 2; i >= 0; i-- {
		getResult(pOne, pathPoint[i])
		str = fmt.Sprintf("  %s", srcFile[pathPoint[i].x])
		result = append(result, str)
		pOne = newpoint(pathPoint[i].x+1, pathPoint[i].y+1)
	}
	pEnd := newpoint(srcLen, dstLen)
	getResult(pOne, pEnd)

	for _, line := range result {
		fmt.Printf("%s\n", line)
	}
}
