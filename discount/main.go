package main

import "fmt"

func main() {
	discount := []int{89000, 99000, 71000}
	per := 0.30
	for _, dk := range discount {
		fmt.Println(int(float64(dk) * (1 - per)))
	}
}
