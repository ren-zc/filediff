package diffV2

import (
	"fmt"
	"io/ioutil"
	// "os"
	"bytes"
	"runtime"
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

// func scanGraph(slen int, dlen int) []*point {
// 	theSame := make([]*point, 0, 100)
// 	p = new(0, 0)
// 	theSame = append(theSame, p)
// 	for x := 0; x < slen; x++ {
// 		for y := 0; y < dlen; y++ {
// 			if srcFile[x] == dstFile[y] {
// 				p = new(x, y)
// 				theSame = append(theSame, p)
// 			}
// 		}
// 	}
// }

var path map[*point][]*point

func init() {
	path = map[*point][]*point{}
}

func scanPath(p *point) []*point {
	shortPath := []*point{}
	xlimit := srcLen
	ylimit := dstLen
	for x, y := p.x+1, p.y+1; x < xlimit && y < ylimit; x, y = x+1, y+1 {
		// fmt.Printf("%d,%d\t", xlimit, ylimit)
		// fmt.Printf("%v\t", newpoint(p.x-1, p.y-1))
		if bytes.Equal(srcFile[x], dstFile[y]) {
			pn := newpoint(x, y)
			shortPath = append(shortPath, pn)
			// fmt.Printf("%v ", pn)
			// shortPath[p] = pchild
			return shortPath
		}
		for i := x + 1; i < xlimit; i++ {
			if bytes.Equal(srcFile[i], dstFile[y]) {
				xlimit = i
				pi := newpoint(i, y)
				// fmt.Printf("%v\t", newpoint(p.x-1, p.y-1))
				// fmt.Printf("%v ", pi)
				// fmt.Printf("%d,%d ", xlimit, ylimit)
				// fmt.Printf("i:%d y:%d\n", i, y)
				shortPath = append(shortPath, pi)
				break
			}
		}
		for j := y + 1; j < ylimit; j++ {
			if bytes.Equal(dstFile[j], srcFile[x]) {
				ylimit = j
				pj := newpoint(x, j)
				// fmt.Printf("%v\t", newpoint(p.x-1, p.y-1))
				// fmt.Printf("%v ", pj)
				// fmt.Printf("%d,%d ", xlimit, ylimit)
				// fmt.Printf("x:%d j:%d\n", x, j)
				shortPath = append(shortPath, pj)
				break
			}
		}
		// fmt.Println()
	}
	return shortPath
}

func getPath(p *point) map[*point][]*point {
	// path :=

	ps := scanPath(p)
	// pp := newpoint(p.x-1, p.y-1)
	if len(ps) == 0 {
		return nil
	}
	path[p] = ps
	for _, pn := range ps {
		// pm := newpoint(pn.x+1, pn.y+1)
		getPath(pn)
	}
	return path
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

func Diff(src string, dst string) {
	srcFile, _ = readFile(src)
	dstFile, _ = readFile(dst)
	srcLen = len(srcFile)
	dstLen = len(dstFile)
	p0 := newpoint(-1, -1)
	getPath(p0)
	// fmt.Println(path)
	for k, v := range path {
		fmt.Printf("%v\t%v\n", k, v)
	}
}

// here
// redo 去重
