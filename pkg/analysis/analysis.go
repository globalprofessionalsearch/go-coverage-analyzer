// Package analysis is responsible for all business logic associated with
// conducting a coverage analysis
package analysis

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	//nolint:exptostd
	"golang.org/x/exp/constraints"
)

var (
	ErrUnsafeCoverProfilePath = errors.New("cover profile path must be prefixed with ./")
)

func Run(coverProfileFileName string, coverageStandard float32) (*ProjectSummary, error) {
	// open the raw coverage file
	// jumping through some security hoops
	// see https://securego.io/docs/rules/g304
	sanitizedPath := filepath.Clean(coverProfileFileName)
	safeBasePath := "./"
	if !strings.HasPrefix(sanitizedPath, safeBasePath) {
		return nil, ErrUnsafeCoverProfilePath
	}
	coverageFile, err := os.Open(sanitizedPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading coverage file")
	}
	//nolint:errcheck
	defer coverageFile.Close()

	// convert each raw line into a block object
	nonDistinctBlocks := []*Block{}
	scanner := bufio.NewScanner(coverageFile)
	firstLine := true
	for scanner.Scan() {
		if firstLine {
			firstLine = false
			continue
		}

		block := new(Block)
		err = block.HydrateFromRawLine(scanner.Text())
		if err != nil {
			return nil, errors.Wrap(err, "unable to read")
		}
		nonDistinctBlocks = append(nonDistinctBlocks, block)
	}

	// The same block can be present multiple times in the raw data
	// (one time for each package tested) we need to merge these values into
	// a single block.
	groupedBlocks := lo.GroupBy(nonDistinctBlocks, func(block *Block) string {
		return block.BlockComponent
	})

	distinctBlocks := []*Block{}
	for _, blocks := range groupedBlocks {
		distinctBlock := new(Block)
		for _, block := range blocks {
			distinctBlock.CallCount += block.CallCount

			// these first values shouldn't change for each block within a
			// block component. just overwriting the value each time is
			// inefficient but it's simple
			distinctBlock.BlockComponent = block.BlockComponent
			distinctBlock.FileName = block.FileName
			distinctBlock.PackageName = block.PackageName
			distinctBlock.RawLineItem = block.RawLineItem
			distinctBlock.StatementCount = block.StatementCount
		}
		distinctBlocks = append(distinctBlocks, distinctBlock)
	}

	blocksByPackage := lo.GroupBy(distinctBlocks, func(block *Block) string {
		return block.PackageName
	})

	projectSummary := new(ProjectSummary)
	packageSummaries := []*PackageSummary{}
	for packageName, blocks := range blocksByPackage {
		blockCount := 0
		blockCallCount := 0
		for _, block := range blocks {
			blockCount++
			if block.CallCount > 0 {
				blockCallCount++
			}
		}
		packageSummaries = append(packageSummaries, &PackageSummary{
			PackageName:    packageName,
			BlockCount:     blockCount,
			BlockCallCount: blockCallCount,
			//nolint:mnd
			CoveragePercentage:    100.0 * safeDivide(blockCallCount, blockCount),
			BlocksNotCoveredCount: blockCount - blockCallCount,
		})

		projectSummary.PackageCount++
		projectSummary.BlockCount += blockCount
		projectSummary.BlockCallCount += blockCallCount
	}

	projectSummary.BlocksNotCoveredCount = projectSummary.BlockCount - projectSummary.BlockCallCount

	//nolint:mnd
	projectSummary.CoveragePercentage = 100 * safeDivide(projectSummary.BlockCallCount, projectSummary.BlockCount)
	projectSummary.CoverageStandard = coverageStandard
	projectSummary.CoverageStandardMet = projectSummary.CoveragePercentage >= projectSummary.CoverageStandard

	sort.SliceStable(packageSummaries, func(i, j int) bool {
		return int(packageSummaries[i].CoveragePercentage) > int(packageSummaries[j].CoveragePercentage)
	})

	projectSummary.PackageSummaries = packageSummaries

	return projectSummary, nil
}

func safeDivide[T constraints.Integer | constraints.Float](a, b T) float32 {
	if b == 0 {
		return 0
	}
	return float32(a) / float32(b)
}
