package main

import (
	"fmt"

	"github.com/brianseitel/chomper"
)

func main() {
	c := chomper.New(3000)

	c.Set(1).Set(3).Set(500)

	fmt.Println(c.Count())
}
