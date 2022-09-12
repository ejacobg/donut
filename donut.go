package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	fmt.Println(int(os.Stdin.Fd()), int(os.Stdout.Fd()))
	fmt.Println(term.GetSize(int(os.Stdin.Fd())))
	fmt.Println(term.GetSize(int(os.Stdout.Fd())))
}
