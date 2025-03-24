package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "go-coverage-analysis",
	Short: "Utilities for the platform api codebase.",
}

func init() {
	RootCmd.AddCommand(
		RunCommand,
	)
}

func main() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
