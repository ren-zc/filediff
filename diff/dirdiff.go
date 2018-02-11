package diff

import (
	"fmt"
)

// DiffOnly used for directions diff of gosync, not means to file diff.
// src: file names and file md5 sums of remote host;
// dst: file names and file md5 sums of local host.
func DiffOnly(src []string, dst []string) ([]string, []string) {
	srcFile = src
	dstFile = dst
	srcLen = len(srcFile)
	dstLen = len(dstFile)

	if dstLen == 0 {
		return src, dst
	}

	pTmp := newpoint(-1, -1)
	getPath(pTmp)
	pathPoint := getMostDepth()

	// rm: src newer than dst or created files.
	rm := []string{}
	// add: files that need update in dst or has been deleted in src.
	add := []string{}
	pOne := newpoint(0, 0)
	getResult := func(pOne, pPoint *point) {
		for j := pOne.x; j < pPoint.x; j++ {
			rm = append(rm, srcFile[j])
		}
		for j := pOne.y; j < pPoint.y; j++ {
			add = append(add, dstFile[j])
		}
	}
	for i := len(pathPoint) - 2; i >= 0; i-- {
		getResult(pOne, pathPoint[i])
		pOne = newpoint(pathPoint[i].x+1, pathPoint[i].y+1)
	}
	pEnd := newpoint(srcLen, dstLen)
	getResult(pOne, pEnd)

	fmt.Println(rm)
	fmt.Println(add)

	return rm, add
}
