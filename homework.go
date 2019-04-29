package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"homework/tools"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func BufWriterFile(ff string, mes string) {
	fp, err := os.OpenFile(ff, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0751)
	if err != nil {
		fmt.Printf("open file %s failed: %v\n", ff, err)
		panic(err)
	}
	defer fp.Close()
	writer := bufio.NewWriter(fp)
	_, err = writer.Write([]byte(mes))
	if err != nil {
		fmt.Printf("write file %s failed: %v\n", ff, err)
	}
	writer.Flush()
}

/*
按offset偏移量读取一行数据, ffp切片维护每个文件的偏移量值
*/
func ReadLine(fp *os.File, offset int64) (string, int64) {
	buffer := make([]byte, 1)
	var res []byte
	for {
		_, err := fp.ReadAt(buffer, offset)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		if buffer[0] == byte('\n') {
			break
		} else {
			res = append(res, buffer[0])
			offset ++
		}
	}
	return string(res), offset
}

func output(fn string, mp map[string]int, ft int) {
	arr := []tools.Dict{}
	for key, val := range mp {
		arr = append(arr, tools.Dict{key, val})
	}
	sort.Sort(tools.Dicts(arr))
	fileStr := fn + "_" + fmt.Sprintf("%04d", ft)
	for _, v := range arr {
		reStr := v.Str + "\t" + fmt.Sprintf("%d", v.Cnt) + "\n"
		BufWriterFile(fileStr, reStr)
	}
}

/*
按topN大小计算优先队列叶子节点的起始下标，这个队列的最小值通过这些叶子节点选出
*/
func calcLeaf(n int) int {
	tmp := n
	for {
		if 2*tmp+1 > n {
			tmp--
			continue
		}
		return tmp + 1
	}
}

/*
按字符串的出现次数维护一个topN容量的优先队列
*/
func heapInsert(heapArr *tools.DictHeap, node tools.Dict, leaf int, n int) {
	heap.Push(heapArr, node)
	if len(*heapArr) > n {
		var min, minid int
		for i := leaf; i <= n; i++ {
			if i == leaf {
				min = (*heapArr)[i].Cnt
				minid = i
			} else {
				if (*heapArr)[i].Cnt < min {
					min = (*heapArr)[i].Cnt
					minid = i
				}
			}
		}
		heap.Remove(heapArr, minid)
	}
}

/*
k路归并各个子文件，第一遍遍历每个文件的第一行找出字典序最小的字符串，然后找出该字符串
总的出现次数，被选中的字符串所在文件游标后移
*/
func merge(fn string, fnum int, topn int) *tools.DictHeap {
	fcur := make([]int64, fnum)
	ffp := make([]*os.File, fnum)
	var minStr string
	var cnt, minFile int
	var minOffset int64
	var err error
	res := &tools.DictHeap{}
	leaf := calcLeaf(topn)
	for i := 0; i < fnum; i++ {
		fileStr := fn + "_" + fmt.Sprintf("%04d", i)
		ffp[i], err = os.OpenFile(fileStr, os.O_RDONLY, 0751)
		if err != nil {
			panic(err)
		}
	}
	defer func() {
		for i := 0; i < fnum; i++ {
			ffp[i].Close()
		}
	}()
	for {
		flag := 0
		for i := 0; i < fnum; i++ {
			tmpStr, offset := ReadLine(ffp[i], fcur[i])
			if tmpStr == "" {
				continue
			}
			s := strings.Split(tmpStr, "\t")
			if flag == 0 || s[0] < minStr {
				minStr = s[0]
				cnt, err = strconv.Atoi(s[1])
				if err != nil {
					panic(err)
				}
				minFile = i
				minOffset = offset + 1
				flag |= 1
			}
		}
		fcur[minFile] = minOffset
		for i := 0; i < fnum; i++ {
			if i == minFile {
				continue
			}
			tmpStr, offset := ReadLine(ffp[i], fcur[i])
			if tmpStr == "" {
				continue
			}
			s := strings.Split(tmpStr, "\t")
			if minStr == s[0] {
				x, err := strconv.Atoi(s[1])
				if err != nil {
					panic(err)
				}
				cnt += x
				fcur[i] = offset + 1
			}
		}
		if flag == 0 {
			break
		}
		heapInsert(res, tools.Dict{minStr, cnt}, leaf, topn)
	}
	return res
}

func main() {
	var fileName string
	var topn, maxMem int
	flag.StringVar(&fileName, "file", "", "the file name you choose")
	flag.IntVar(&topn, "topn", 10, "output the top n numbers")
	flag.IntVar(&maxMem, "maxmem", 1024*1024*1024, "the max memory size you can use")
	flag.Parse()
	// 记录字节数和文件个数
	var bcnt, fcnt int
	mmp := make(map[string]int)
	fp, err := os.OpenFile(fileName, os.O_RDONLY, 0751)
	if err != nil {
		fmt.Printf("open file %s failed: %v\n", fileName, err)
		panic(err)
	}
	reader := bufio.NewReader(fp)
	var rstr string
	for {
		rstr, err = reader.ReadString('\n')
		if err == nil {
			str := strings.TrimSuffix(rstr, "\n")
			if _, ok := mmp[str]; ok {
				mmp[str]++
				continue
			}
			if bcnt + len(str) + 4 > maxMem {
				output(fileName, mmp, fcnt)
				bcnt = 0
				fcnt++
				mmp = make(map[string]int)
			}
			bcnt += (len(str) + 4)
			mmp[str] = 1
		} else if err == io.EOF {
			if bcnt != 0 {
				output(fileName, mmp, fcnt)
				fcnt++
			}
			break
		} else {
			fmt.Printf("read file %s failed: %v\n", fileName, err)
			panic(err)
		}
	}
	fp.Close()
	fmt.Println("dispatch finishied")
	result := merge(fileName, fcnt, topn)
	for {
		if len(*result) <= 0 {
			break
		}
		fmt.Println(heap.Pop(result))
	}
}
