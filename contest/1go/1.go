package main

import ( 
	"fmt" 
	"math" 
)

func main() { 
	var n int 
	fmt.Scan(&n)

	shurn := 1
	shursum := 0
	res := 0

	for i := 0; i < n; i++ {
		shurn = 2*(i+1) - 1
		shursum += int(math.Pow(float64(shurn), 2))
		res = int(math.Pow(float64(shurn), 3)) - shursum
		fmt.Print(res, " ")
	}
}