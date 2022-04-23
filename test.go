package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now().YearDay()
	fmt.Println(t)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
