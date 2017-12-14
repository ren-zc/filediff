// By jacenr
// Create: 2017-12-10
/*
 * Usage:
 * import "github.com/jacenr/filediff/diff"
 * result, _ := diff.Diff("file1", "file2")
 */
package diff

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var srcFile []string
var dstFile []string
var srcLen int
var dstLen int
var path = map[*point][]*point{}
var newed = map[string]*point{}

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

// Check the point whether is exist.
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

// Get all shortcut paths of a point.
func scanPath(p *point) []*point {
	shortPath := []*point{}
	xlimit := srcLen
	ylimit := dstLen
	for x, y := p.x+1, p.y+1; x < xlimit && y < ylimit; x, y = x+1, y+1 {
		if srcFile[x] == dstFile[y] {
			pn := checkNew(x, y)
			shortPath = append(shortPath, pn)
			return shortPath
		}
		for i := x + 1; i < xlimit; i++ {
			if srcFile[i] == dstFile[y] {
				xlimit = i
				pi := checkNew(i, y)
				shortPath = append(shortPath, pi)
				break
			}
		}
		for j := y + 1; j < ylimit; j++ {
			if srcFile[x] == dstFile[j] {
				ylimit = j
				pj := checkNew(x, j)
				shortPath = append(shortPath, pj)
				break
			}
		}
	}
	return shortPath
}

// Put all shortcut paths into a map.
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

// Get the best path.
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

// Read file text.
func readFile(file string) ([]string, error) {
	fileContens := []string{}
	f, FErr := os.Open(file)
	if FErr != nil {
		return nil, FErr
	}
	ScannerF := bufio.NewScanner(f)
	ScannerF.Split(bufio.ScanLines)
	for ScannerF.Scan() {
		fileContens = append(fileContens, ScannerF.Text())
	}
	return fileContens, nil
}

// Output difference of files.
func Diff(src string, dst string) ([]string, error) {
	var fileErr error
	srcFile, fileErr = readFile(src)
	if fileErr != nil {
		return nil, fileErr
	}
	dstFile, fileErr = readFile(dst)
	if fileErr != nil {
		return nil, fileErr
	}
	srcLen = len(srcFile)
	dstLen = len(dstFile)
	pTmp := newpoint(-1, -1)
	getPath(pTmp)
	// for k, v := range path {      **
	// 	fmt.Printf("%v\t%v\n", k, v) ** FOR DEBUG
	// }                             **
	pathPoint := getMostDepth(pTmp)
	// fmt.Println(pathPoint)        ** FOR DEBUG

	// output
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
	return result, nil
}
