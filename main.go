package main

import (
	"fmt"
	"os"

	"github.com/inari111/hatebu/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
