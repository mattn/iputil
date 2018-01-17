package main

import (
	"fmt"
	"os"

	"github.com/mattn/iputil"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage of %s: [QUERY]\n", os.Args[0])
		os.Exit(1)
	}
	rngs, err := iputil.Ranges(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", err)
		os.Exit(1)
	}
	for _, r := range rngs {
		fmt.Println(r.CIDR())
	}
}
