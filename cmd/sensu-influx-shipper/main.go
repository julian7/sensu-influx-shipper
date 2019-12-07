package main

import (
	"fmt"
	"os"
)

var version = "SNAPSHOT"

func main() {
	rt := NewRuntime()

	cmd, err := rt.rootCmd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		_ = rt.Log("error", err)

		os.Exit(1)
	}
}
