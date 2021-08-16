package main

import "fmt"

const name = "123"

var func01 = func(int){}

func main() {
	go func(){
		fmt.Println("World")
	}()
	fmt.Println("Hello")
}
