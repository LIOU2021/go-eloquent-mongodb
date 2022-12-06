package main

import "fmt"

func main() {
	defer fmt.Println("33333")
	fmt.Println(13/6 == 2)
	fmt.Println(17 % 3)
	echo()
}

func echo() {
	defer fmt.Println("2222")
	defer fmt.Println("2222-111")
	fmt.Println("1111")
}
