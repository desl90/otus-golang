package main

import (
	"fmt"

	"github.com/desl90/otus-golang/hw02_unpack_string/unpacker"
)

func main() {
	str := `qwe\\\3`

	unpackedString, err := hw02unpackstring.Unpack(str)

	if err != nil {
		fmt.Printf("Exception: %v", err)

		return
	}

	fmt.Printf("Unpacked string: %v\n", unpackedString)
}
