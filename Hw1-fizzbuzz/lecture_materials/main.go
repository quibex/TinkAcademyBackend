package main

import (
	"fmt"
	"unsafe"
)

func stringExample1() {
	str := "世界世界世界"

	// 18 = 6 символов * 3 байта на символ
	fmt.Println(len(str))

	for _, ch := range str {
		fmt.Println(ch, string(ch))
	}

	for i := range []byte(str) {
		fmt.Printf("%b \n", str[i])
	}
}

func stringExampleAscii() {
	str := "abc"

	// 3 = 3 символов * 1 байта на символ
	fmt.Println(len(str))

	for _, ch := range str {
		fmt.Println(ch, string(ch))
	}

	for i := range []byte(str) {
		fmt.Printf("%b \n", str[i])
	}
}

func typingExample() {
	fmt.Println(max(1, 2))
	fmt.Println(max(12.1, 123.1))
}

func arrayExmaple() {
	arr := []int{1, 1, 2, 4}

	fmt.Println(unsafe.Sizeof(arr), len(arr), cap(arr)) // 24 4 4

	arr2 := make([]int, 100)
	fmt.Println(unsafe.Sizeof(arr2), len(arr2), cap(arr2)) // 24 100 100
	fmt.Println(arr2[50])

	arr3 := make([]int, 0, 100)
	fmt.Println(unsafe.Sizeof(arr3), len(arr3), cap(arr3)) // 24 0 100
	fmt.Println(arr3[50])                                  // panic: runtime error: index out of range [50] with length 0
}

func arrayAgainExample() {
	//type slice struct {
	//	array unsafe.Pointer
	//	len   int
	//	cap   int
	//}

	arr := []int{1, 1, 2, 4}

	arrNew := arr[1:]

	arr[2] = 1000

	fmt.Println(arr, len(arr), cap(arr))          // [1 1 1000 4] 4 4
	fmt.Println(arrNew, len(arrNew), cap(arrNew)) // [1 1000 4] 3 3

	arrNew = append(arrNew, []int{99, 9, 8}...)

	arr[2] = 2000

	fmt.Println(arr, len(arr), cap(arr))          // [1 1 2000 4] 4 4
	fmt.Println(arrNew, len(arrNew), cap(arrNew)) // [1 1000 4 99 9 8] 6 6
}

func main() {
	// SELECT a, b, c FROM table LIMIT 1000;
	// []int -> []string
	//a := make([]int, 1000)

	x := []int{1, 2, 3, 4, 5, 6, 7}
	a := x[2:] // 2 3 4 5 6
	b := x[:3] // 0 1 2

	func(arr []int) { // передача всегда по ЗНАЧЕНИЮ
		fmt.Printf("%p %v \n", arr, arr) // 0xc0000200c0
		fmt.Println(len(arr), cap(arr))  // 3 7

		arr = append(arr, 100)

		fmt.Printf("%p %v \n", arr, arr) // 0xc0000200c0
		fmt.Println(len(arr), cap(arr))  // 4 7
	}(b)
	//b = append(b, 100)
	fmt.Printf("%p %v \n", b, b) // 0xc0000200c0
	fmt.Println(len(b), cap(b))  // 3 7

	fmt.Println(len(a), cap(a)) // 5 5
	fmt.Println(len(b), cap(b)) // 4 7
	fmt.Println(x)              // [1 2 3 100 5 6 7]
	fmt.Println(a)              // [3 100 5 6 7]
	fmt.Println(b[:5])          // [1 2 3 100]
	fmt.Println(a[:6])          // slice bounds out of range [:6] with capacity 5

	fmt.Println(a[len(a)-1])
}
