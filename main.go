package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "importpaths":
		paths, err := importpaths(flag.Arg(1))

		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}

		for _, p := range paths {
			fmt.Println(p)
		}
	}
}
