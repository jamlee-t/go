package data

import "fmt"

const name = "123"

func Hello() {
	go func(){
		fmt.Println("World")
	}()
	fmt.Println("Hello")
}
