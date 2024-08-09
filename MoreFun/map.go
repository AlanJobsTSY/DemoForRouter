package main

import "log"

var test = map[int]string{1: "1", 2: "2", 3: "3"}
var cnt = map[int]int{1: 0, 2: 0, 3: 0}

func main() {
	for i := 0; i < 1000; i++ {
		for key, _ := range test {
			cnt[key] += 1
			break
		}
	}
	log.Printf("%d %d %d", cnt[1], cnt[2], cnt[3])

}
