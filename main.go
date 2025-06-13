package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(rand.Intn(100)) // always same result every time you run the program unless you seed
}
