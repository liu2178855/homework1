package tools

import(
)

type Dict struct {
	Str string
	Cnt int
}

type Dicts []Dict

func (d Dicts) Len() int { return len(d) }

func (d Dicts) Less(i, j int) bool {
        return d[i].Str < d[j].Str
}

func (d Dicts) Swap(i, j int) {
        d[i], d[j] = d[j], d[i]
}

type DictHeap []Dict

func (dh DictHeap) Len() int { return len(dh) }

func (dh DictHeap) Swap(i, j int) {
	dh[i], dh[j] = dh[j], dh[i]
}

func (dh DictHeap) Less(i, j int) bool {
	return dh[i].Cnt > dh[j].Cnt
}

func (dh *DictHeap) Push(d interface{}) {
	*dh = append(*dh, d.(Dict))
}

func (dh *DictHeap) Pop() interface{} {
	n := len(*dh)
	x := (*dh)[n-1]
	*dh = (*dh)[:n-1]
	return x
}


