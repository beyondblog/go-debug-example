package main

import (
	"fmt"
	"github.com/beyondblog/go-debug-example/lib"
	"os"
	"runtime"
	//"runtime/debug"
)

func main() {
	var modify string
	argsLen := len(os.Args)
	if argsLen < 2 {
		fmt.Printf("Usage go-debug-example [username] \r\n")
		os.Exit(-1)
	}
	username := os.Args[1]
	var password string
	fmt.Printf("%s welcome!\r\nplease input password:", username)
	fmt.Scanf("%s", &password)
	fmt.Printf("%s  password: %s\r\n", username, password)

	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
		if i == 5 {
			modify = "modify!"
		}
	}

	fmt.Println(lib.Add(sum, 10))

	runtime.Breakpoint()
	//debug.PrintStack()
	fmt.Println(sum)
	fmt.Println(modify)
}
