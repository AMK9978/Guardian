package main

import (
	"fmt"
	"os"

	"guardian/cmd/guardian"
)

func main() {
	if err := guardian.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
