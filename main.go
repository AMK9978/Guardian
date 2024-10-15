package main

import (
	"fmt"
	"guardian/cmd/guardian"
	"os"
)

func main() {
	if err := guardian.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
