package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

func main() {
	// 创建一个 TreeMap
	m := treemap.NewWithIntComparator()

	// 插入键值对
	m.Put(1, "one")
	m.Put(3, "three")
	m.Put(5, "five")
	m.Put(1848953500, "seven")
	m.Remove(3)
	// 查找不小于 4 的第一个键

	k, v := m.Ceiling(2159219264)

	fmt.Printf("%d,%s", k, v)
	fmt.Println(m.Size())
}
