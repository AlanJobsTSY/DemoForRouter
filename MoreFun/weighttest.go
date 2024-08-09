package main

import (
	"log"
)

var a = []int{2, 3, 5}
var c = []int{3, 1, 2}

func main() {
	for i := 1; i <= 20; i++ {
		minn := 0
		index := -1
		for j := 0; j < len(a); j++ {
			c[j] += a[j]
		}
		for j := 0; j < len(a); j++ {
			if c[j] > minn {
				index = j
			}
		}
		c[index] -= 10
		log.Printf("%d", index)
	}
}
