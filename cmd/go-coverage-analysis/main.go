package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-coverage-analysis",
	Short: "Utilities for the platform api codebase.",
}

func init() {
	rootCmd.AddCommand(
		runCommand,
	)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
