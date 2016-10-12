package main

import (
	"fmt"
	"time"
)

/////////////////////// simple statement //////////////////////////////////////
func state_test() {
	for j := 1; j <= 9; j++ {
		fmt.Println(j)
	}

	// while loop
	for {
		fmt.Println("loop")
		break
	}

	// if state
	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}

	//switch
	// You can use commas to separate multiple expressions in the same case
	// statement. We use the optional default case in this example as well.
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("it's the weekend")
	default:
		fmt.Println("it's a weekday")
	}
	// switch without an expression is an alternate way to express if/else logic.
	// Here we also show how the case expressions can be non-constants.
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("it's before noon")
	default:
		fmt.Println("it's after noon")
	}

	// array
	// Use this syntax to declare and initialize an array in one line.
	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("dcl:", b)
	// Array types are one-dimensional, but you can compose types to build
	// multi-dimensional data structures.
	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)

	// Unlike arrays, slices are typed only by the elements they contain
	s := make([]string, 3)
	fmt.Println("emp:", s)
	s = append(s, "hello")
	fmt.Println("append:", s)
	//We can declare and initialize a variable for slice in a single line as well.
	ts := []string{"g", "h", "i"}
	fmt.Println("dcl:", ts)
	//Slices can be composed into multi-dimensional data structures. The length
	//of the inner slices can vary, unlike with multi-dimensional arrays.
	twoDS := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoDS[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoDS[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoDS)

	// map
	var m = make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	delete(m, "a")
	n := map[string]int{"k1": 10, "k2": 20}
	fmt.Println("Map, m=", m, ", n=", n)

	// range
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}
}

//////////////////// Variadic Functions //////////////////////////////////////
// fmt.Println(sum(1, 3, 4))
// nums := []int{1, 2, 3, 4}
// fmt.Println(sum(nums...))
func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

//////////////////// Closures ////////////////////////////////////////////////
// nextInt := intSeq()
func intSeq() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

/////////////////// struct and methods ////////////////////////////////////////
type rect struct {
	width, height int
}

// Go automatically handles conversion between values and pointers for method
// calls. You may want to use a pointer receiver type to avoid copying on method
// calls or to allow the method to mutate the receiving struct.
func (r rect) area() int {
	return r.width * r.height
}
func (r rect) perim() int {
	return 2*r.width + 2*r.height
}

///////////////////// Interface ///////////////////////////////////////////////
type geometry interface {
	area() int
	perim() int
}

func main() {
	state_test()
	// interace
	var geo geometry
	geo = rect{width: 2, height: 3}
	fmt.Println(geo.area())
}
