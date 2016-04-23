package main

import (
	"fmt"

	"github.com/hearts.zhang/mecab"
)

func main() {
	mc, _ := mecab.New()
	if mc != nil {
		defer mc.Destroy()
	}
	var x = "北京大学香港中文大学西窗释放"
	nodes, _ := mc.Sparse2(x)
	for _, node := range nodes {
		fmt.Println(node)
	}
}
