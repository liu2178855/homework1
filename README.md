## 安装

	go build main.go
## 测试
	go test -v main.go homework_test.go

## 使用

	./main --file=testdata --maxmem=524288 --topn=10 

## 思路

	1. 读入源文件存入map，当map大于maxmem时，按字符串字典序排序后输出到文件，每一行记录字符串及出现次数
	2. 归并各个子文件字符串出现次数，k路归并，然后把最终的字符串及其次数写入节点个数为topn的堆（优先队列）中
	3. 当队列超过topn时，remove掉所有叶子节点中最小的那个，然后按顺序输出优先队列
