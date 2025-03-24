package main

import (
	"fmt"
	"os"

	"github.com/globalprofessionalsearch/go-coverage-analyzer/pkg/analysis"
	"github.com/spf13/cobra"
)

var (
	coverProfileFileName    string
	coverageStandardPercent float32
)

func init() {
	RunCommand.Flags().StringVarP(&coverProfileFileName, "coverprofile", "c", "coverage.out", "the name of the coverage profile file to analyze")
	RunCommand.Flags().Float32VarP(&coverageStandardPercent, "standard", "s", 75, "minimum coverage percentage required (exit with error if not met)")
}

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: "Conducts codebase coverage analysis.",
	Long: `Description: 
	Condenses coverage file output so that each code-block is only represented a
	single time. This is helpful when tests have been run in coverpkg mode. Once
	the coverage output is consolidated, coverage statistics are compiled by
	package as well as across the entire codebase and returned via stdout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("\nChecking coverage...\n")

		projectSummary, err := analysis.Run(*&coverProfileFileName, *&coverageStandardPercent)
		if err != nil {
			return err
		}

		fmt.Println("\nPackageName", "CoveragePercentage", "Blocks", "BlocksCovered", "CoveragePercentage")
		for _, packageSummary := range projectSummary.PackageSummaries {
			blocksCovered := packageSummary.BlockCount - packageSummary.BlocksNotCoveredCount
			fmt.Printf("%s\t%d\t%d\t%.2f\n", packageSummary.PackageName, packageSummary.BlockCount, blocksCovered, packageSummary.CoveragePercentage)
		}

		fmt.Printf("\nProject Summary\n")
		fmt.Println("PackageCount", "CoveragePercentage")
		fmt.Printf("%d\t%.2f%%\n", projectSummary.PackageCount, projectSummary.CoveragePercentage)

		fmt.Printf("\nChecking coverage against standard...\n")
		fmt.Printf("Actual: %.2f\n", projectSummary.CoveragePercentage)
		fmt.Printf("Standard: %.2f\n", projectSummary.CoverageStandard)

		if projectSummary.CoverageStandardMet {
			fmt.Printf("\nCoverage standard met!\n")
		} else {
			fmt.Printf("\nCoverage percentage doesn't meet standard.\n")

			// We exit 1 instead of returning an error. Returning an error
			// would cause cobra to display the help message.
			os.Exit(1)
		}

		return nil
	},
}
