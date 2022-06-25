package main

import (
	"fmt"
	"os"
)

func main() {
	f, _ := os.Getwd()
	fmt.Println(f)
}
