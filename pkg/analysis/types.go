package analysis

import (
	"path"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type Block struct {
	// RawLineItem is the raw coverage line that represents a single analyzed
	// code-block where a code-block is a group of statements.
	RawLineItem string

	// BlockComponent is the first slug of the raw line that defines the
	// package, file, start line, start column, end line, and end column.
	BlockComponent string

	// PackageName is the name of the package that was analyzed.
	PackageName string

	// FileName is the name of the specific file that was analyzed.
	FileName string

	// Statement count is the name of statements in a particular block of code.
	StatementCount int

	// CallCount is the number of times that block of code was called.
	CallCount int
}

type PackageSummary struct {
	// PackageName is the name of the package being summarized
	PackageName string

	// BlockCount is the number of statement blocks present in the package.
	BlockCount int

	// BlockCallCount is the number of blocks that were called at least once.
	BlockCallCount int

	// CoveragePercentage is percentage of blocks that were called at least
	// once (expressed as a value from 0 to 100).
	CoveragePercentage float32

	// BlocksNotCoveredCount is the number of blocks that did not receive any
	// calls.
	BlocksNotCoveredCount int
}

type ProjectSummary struct {
	// PackageCount is the total number of packages
	PackageCount int

	// BlockCount is the number of statement blocks present in the project.
	BlockCount int

	// BlockCallCount is the number of blocks that were called at least once.
	BlockCallCount int

	// BlocksNotCoveredCount is the number of blocks that did not receive any
	// calls.
	BlocksNotCoveredCount int

	// CoveragePercentage is percentage of blocks that were called at least
	// once (expressed as a value from 0 to 100).
	CoveragePercentage float32

	// CoverageStandard is the coverage-standard percentage.
	CoverageStandard float32

	// CoverageStandardMet is true if the standard is met else false.
	CoverageStandardMet bool

	// PackageSummaries are the summaries for each analyzed package.
	PackageSummaries []*PackageSummary
}

func (c *Block) HydrateFromRawLine(line string) error {
	c.RawLineItem = line

	firstSplit := strings.Split(line, " ")
	c.BlockComponent = firstSplit[0]
	c.StatementCount = lo.Must(strconv.Atoi(firstSplit[1]))
	c.CallCount = lo.Must(strconv.Atoi(firstSplit[2]))

	secondSplit := strings.Split(firstSplit[0], ":")
	c.FileName = secondSplit[0]
	c.PackageName = path.Dir(c.FileName)
	return nil
}
