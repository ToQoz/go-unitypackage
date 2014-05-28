package main

import (
	"fmt"
	"os"
)

func uninstall(unitypackagePath string) error {
	paths, err := list(unitypackagePath)

	if err != nil {
		return err
	}

	for _, p := range paths {
		_, err := os.Stat(p)

		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			continue
		}

		err = os.Remove(p)

		if err != nil {
			return err
		}
	}

	return nil
}
