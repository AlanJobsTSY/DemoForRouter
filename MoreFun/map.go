package main

import (
	"log"
	"time"
)

var test = map[int]string{1: "1", 2: "2", 3: "3"}
var cnt = map[int]int{1: 0, 2: 0, 3: 0}

func main() {

	for {
		log.Printf("1")
		time.Sleep(5 * time.Second)
	}
}
