package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	var z0, z1 float64
	z0, z1 = 0, 10000
	for math.Abs(z1-z0) > 0.001 {
		z0 = z1
		z1 = z1 - (z1*z1-x)/(2*z1)
	}
	return z1
}

func main() {
	fmt.Println(Sqrt(2))
}
