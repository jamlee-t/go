// pos， 注释测试，对pos是否有影响
// 两行注释呢
package data
import (
	"time"
)
import "fmt"

// 常量Smt
const name = 123

// 声明变量 Smt
var func01 = func(int){}

// 函数声明 Smt
func Hello() {
	time.Now()
	// 内部包含 Smt
	go func(){
		fmt.Println("World")
	}()
	fmt.Println("Hello")
	fmt.Println(name)
}
