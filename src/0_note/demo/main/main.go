package main

import "fmt"
import (
	"time"
)

const name = "123"

var func01 = func(int){}

func main() {
	time.Now()
	go func(){
		fmt.Println("World")
	}()
	fmt.Println("Hello")
}
