package diff

import (
	"fmt"
	// "log"
	"math/rand"
	"os"
	"testing"
	"time"
)

// Diff的精简形式返回最大相同行数，仅用于测试。
func Diff2(src string, dst string) (int, error) {
	var fileErr error
	srcFile, fileErr = readFile(src)
	if fileErr != nil {
		return 0, fileErr
	}
	dstFile, fileErr = readFile(dst)
	if fileErr != nil {
		return 0, fileErr
	}
	srcLen = len(srcFile)
	dstLen = len(dstFile)
	pTmp := newpoint(-1, -1)
	getPath(pTmp)
	pathPoint := getMostDepth()
	return len(pathPoint), nil
}

// 创建随机内容文件, 返回文件名和已知的相同行数
func CreateRandFile(t *testing.T) (string, string, int) {
	rand.Seed(time.Now().UnixNano())
	var theSameRandLen int
	var srcRandLen int
	var dstRandLen int
	theSameRandLen = rand.Intn(1000)
	for theSameRandLen < 500 {
		theSameRandLen = rand.Intn(1000)
	}
	for srcRandLen < 500 {
		srcRandLen = rand.Intn(1000)
	}
	for dstRandLen < 500 {
		dstRandLen = rand.Intn(1000)
	}
	// theSame := []string{}
	theSame := make([]string, 0, theSameRandLen)
	SameStr := "The Same "
	for i := 0; i < theSameRandLen; i++ {
		theSame = append(theSame, SameStr+string(rand.Intn(95)+33)+"\n")
	}
	// srcOnly := []string{}
	srcOnly := make([]string, 0, srcRandLen)
	srcStr := "The Only Src "
	for i := 0; i < srcRandLen; i++ {
		srcOnly = append(srcOnly, srcStr+string(rand.Intn(95)+33)+"\n")
	}
	// dstOnly := []string{}
	dstOnly := make([]string, 0, dstRandLen)
	dstStr := "The Only Dst "
	for i := 0; i < dstRandLen; i++ {
		dstOnly = append(dstOnly, dstStr+string(rand.Intn(95)+33)+"\n")
	}

	now := time.Now()
	fileNameStr := fmt.Sprintf("%d%d%d", now.Day(), now.Hour(), now.Second())

	srcFileName := "../testFile/" + "src" + fileNameStr
	src, srcErr := os.Create(srcFileName)
	if srcErr != nil {
		// log.Fatalln(srcErr)
		t.Fatal(srcErr)
	}
	// theSameDst := []string{}
	theSameDst := make([]string, 0, theSameRandLen)
	for srcOnly != nil || theSame != nil {
		if rand.Intn(2) == 0 && srcOnly != nil {
			src.WriteString(srcOnly[0])
			if len(srcOnly) > 1 {
				srcOnly = srcOnly[1:]
			} else {
				srcOnly = nil
			}
			continue
		}
		if theSame != nil {
			src.WriteString(theSame[0])
			theSameDst = append(theSameDst, theSame[0])
			if len(theSame) > 1 {
				theSame = theSame[1:]
			} else {
				theSame = nil
			}
		}
	}
	src.Close()
	dstFileName := "../testFile/" + "dst" + fileNameStr
	dst, dstErr := os.Create(dstFileName)
	if dstErr != nil {
		// log.Fatalln(dstErr)
		t.Fatal(dstErr)
	}
	for dstOnly != nil || theSameDst != nil {
		if rand.Intn(2) == 0 && dstOnly != nil {
			dst.WriteString(dstOnly[0])
			if len(dstOnly) > 1 {
				dstOnly = dstOnly[1:]
			} else {
				dstOnly = nil
			}
			continue
		}
		if theSameDst != nil {
			dst.WriteString(theSameDst[0])
			if len(theSameDst) > 1 {
				theSameDst = theSameDst[1:]
			} else {
				theSameDst = nil
			}
		}
	}
	dst.Close()
	return srcFileName, dstFileName, theSameRandLen
}

func TestDiff2(t *testing.T) {
	src, dst, length := CreateRandFile(t)
	rlength, DiffErr := Diff2(src, dst)
	if rlength != length+1 || DiffErr != nil { // 加"1"是因为求取rlength的slice多一个元素(-1,-1)。
		t.Error("False")
	}
	var rmErr error
	rmErr = os.Remove(src)
	if rmErr != nil {
		fmt.Println(rmErr)
	}
	rmErr = os.Remove(dst)
	if rmErr != nil {
		fmt.Println(rmErr)
	}
}
