package data

import "fmt"

const name = "123"

var func01 = func(int){}

func Hello() {
	go func(){
		fmt.Println("World")
	}()
	fmt.Println("Hello")
}
