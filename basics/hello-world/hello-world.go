// Our first program will print the classic "hello world"
// message. Here's the full source code.
package main

import "fmt"
import "os"

func main() {
	fmt.Println("hello world")
	fmt.Println(os.Environ())
}
