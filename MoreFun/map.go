package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var test = map[int]string{1: "1", 2: "2", 3: "3"}
var cnt = map[int]int{1: 0, 2: 0, 3: 0}

func main() {

	// 创建 CPU 分析文件
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	// 开始 CPU 分析
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// 模拟一些工作负载
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
	}

	t := time.Now()
	log.Printf("%v", t)
}
