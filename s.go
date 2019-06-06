package main
import (
	"fmt"
)

func main() {
	b := make([]byte, 10)
	c := make([]byte, 4)
	b = append(b, c...)
	fmt.Println(b)
}
