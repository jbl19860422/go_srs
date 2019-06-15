package main

import(
	"fmt"
)

type B struct {
	data int32
}
type A struct {
	data B
}

func main() {
	a := A{}
	fmt.Println(a)
}
