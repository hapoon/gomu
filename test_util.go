package gomu

import "fmt"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func ignoreError(err error) {
	if err != nil {
		fmt.Printf("[WARN]%v\n", err)
	}
}
