package main

import (
	"fmt"
	"github.com/linxlib/charls"
	"os"
)

func main() {
	f, _ := os.Open("lena8b.jls")
	fmt.Print(charls.GetVersion())
	charls.Decode(f)
}
