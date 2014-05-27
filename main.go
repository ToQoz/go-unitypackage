package main

import (
	"flag"
	"fmt"
	"github.com/ToQoz/go-unitypackage/importpath"
	"os"
)

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "importpaths":
		paths, err := importpath.All(flag.Arg(1))

		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}

		for _, p := range paths {
			fmt.Println(p)
		}
	}
}
