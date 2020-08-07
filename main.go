package main

import (
	"fmt"

	"./mosquittoscope"
)

func main() {
	fmt.Println("Hi")
	s := mosquittoscope.NewSettings("boop.yaml")
	fmt.Printf("%q\n", s)
}
