package main

import (
	"testing"
	"os"
	"homework/tools"
	"container/heap"
)

func Test_ReadLine(t *testing.T) {
	fp, err := os.OpenFile("haha", os.O_RDONLY, 0751)
	if err != nil {
		t.Error("按行读取打开失败")
	}
	str, offset := ReadLine(fp, 0)
	if str == "dadbaskbdakd" && offset == 12 {
		t.Log("按行读取测试通过")
	} else {
		t.Error("按行读取失败")
	}
}

func Test_calcLeaf(t *testing.T) {
	x := calcLeaf(5)
	if x == 3 {
		t.Log("计算叶子节点通过")
	} else {
		t.Log("计算叶子节点失败")
	}
}

func Test_heapInsert(t *testing.T) {
	res := &tools.DictHeap{}
	leaf := calcLeaf(5)
	heapInsert(res, tools.Dict{"abc", 1 }, leaf, 5)
	heapInsert(res, tools.Dict{"abcd", 2 }, leaf, 5)
	heapInsert(res, tools.Dict{"bcde", 3 }, leaf, 5)
	heapInsert(res, tools.Dict{"bcdef", 4 }, leaf, 5)
	heapInsert(res, tools.Dict{"bbcde", 5 }, leaf, 5)
	heapInsert(res, tools.Dict{"cbcde", 6 }, leaf, 5)
	if len(*res) == 5 && (*res)[0].Str == "cbcde" {
		t.Log("指定大小优先队列测试通过")
	} else {
		t.Error("指定大小优先队列测试失败")
	}
}

func Test_merge(t *testing.T) {
	res := merge("test_data", 2, 3)
	if len(*res) != 3 {
		t.Error("k路合并求topn数量字符串不通过")
		return
	}
	var tt, tmp tools.Dict
	tt = (heap.Pop(res)).(tools.Dict)
	tmp = tools.Dict{"ddd", 8}
	if tt != tmp {
		t.Error("k路合并求topn数量字符串不通过")
		return
	}
	tt = (heap.Pop(res)).(tools.Dict)
	tmp = tools.Dict{"eee", 6}
	if tt != tmp {
                t.Error("k路合并求topn数量字符串不通过")
                return
        }
	tt = (heap.Pop(res)).(tools.Dict)
	tmp = tools.Dict{"bbb", 3}
	if tt != tmp {
                t.Error("k路合并求topn数量字符串不通过")
                return
        }
	t.Log("k路合并求topn数量字符串通过")
}
