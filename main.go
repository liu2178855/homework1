package main

import (
	"bufio"
	"flag"
	"fmt"
	"homework/tools"
	"strconv"
	"io"
	"os"
	"sort"
	"strings"
	"container/heap"
)

var fileName string
var topn, maxMem int

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
		fmt.Println("write file %s failed: %v\n", ff, err)
	}
	writer.Flush()
}

//这里需要通过offset来改进下性能, 现在这种写法性能太低
func ReadLine(ff string, lineNumber int) string {
	file, _ := os.OpenFile(ff, os.O_RDONLY, 0751)
	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	for fileScanner.Scan() {
		if lineCount == lineNumber {
			return fileScanner.Text()
		}
		lineCount++
	}
	defer file.Close()
	return ""
}

func output(mp map[string]int, ft int) {
	arr := []tools.Dict{}
	for key, val := range mp {
		arr = append(arr, tools.Dict{key, val})
	}
	sort.Sort(tools.Dicts(arr))
	fileStr := fileName + "_" + fmt.Sprintf("%04d", ft)
	for _, v := range arr {
		reStr := v.Str + "\t" + fmt.Sprintf("%d", v.Cnt) + "\n"
		BufWriterFile(fileStr, reStr)
	}
}

func merge(fnum int) {
	fcur := make([]int, fnum)
	var minStr string
	var cnt, minFile int
	var err error
	var leaf int
	res := &tools.DictHeap{}
	temp := topn
	for  {
		if 2 * temp + 1 > topn {
			temp --
			continue
		}
		leaf = temp + 1
		break
	}
	for {
		flag := 0
		for i := 0; i < fnum; i ++ {
			fileStr := fileName + "_" + fmt.Sprintf("%04d", i)
			tmp := strings.TrimSuffix(ReadLine(fileStr, fcur[i] + 1), "\n")
			if tmp == "" {
				continue
			}
			s := strings.Split(tmp, "\t")
			if flag == 0 || s[0] < minStr {
				minStr = s[0]
                                cnt, err = strconv.Atoi(s[1])
				if err != nil {
					panic(err)
				}
				minFile = i
				flag |= 1
			}
		}
		fcur[minFile] ++
		for i := 0; i < fnum; i ++ {
			if i == minFile {
				continue
			}
			fileStr := fileName + "_" + fmt.Sprintf("%04d", i)
                        tmp := strings.TrimSuffix(ReadLine(fileStr, fcur[i] + 1), "\n")
                        if tmp == "" {
                                continue
                        }
			s := strings.Split(tmp, "\t")
			if minStr == s[0] {
				x, err := strconv.Atoi(s[1])
				if err != nil {
					panic(err)
				}
				cnt += x
				fcur[i] ++
			}
		}
		if flag == 0 {
			break
		}
		heap.Push(res, tools.Dict{minStr, cnt})
		if len(*res) > topn {
			var min, minid int
			for i := leaf; i <= topn; i++ {
				if i == leaf {
					min = (*res)[i].Cnt
					minid = i
				} else {
					if (*res)[i].Cnt < min {
						min = (*res)[i].Cnt
						minid = i
					}
				}
			}
			heap.Remove(res, minid)
		}
	}
	for {
		if len(*res) <= 0 {
			break
		}
		fmt.Println(heap.Pop(res))
	}
}

func main() {
	flag.StringVar(&fileName, "file", "", "the file name you choose")
	flag.IntVar(&topn, "topn", 10, "output the top n numbers")
	flag.IntVar(&maxMem, "maxmem", 1024*1024*1024, "the max memory size you can use")
	flag.Parse()
	// byte count and file count
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
			if bcnt+len(str)+4 > maxMem {
				output(mmp, fcnt)
				bcnt = 0
				fcnt++
				mmp = make(map[string]int)
			}
			bcnt += (len(str) + 4)
			mmp[str] = 1
		} else if err == io.EOF {
			if bcnt != 0 {
				output(mmp, fcnt)
				fcnt ++
			}
			break
		} else {
			fmt.Printf("read file %s failed: %v\n", fileName, err)
			panic(err)
		}
	}
	fmt.Println("dispatch finishied")
	merge(fcnt)
	fp.Close()
}
