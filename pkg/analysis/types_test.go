package analysis

import (
	"testing"
)

func TestBlock_HydrateFromRawLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		wantFileName  string
		wantPackage   string
		wantStatCount int
		wantCallCount int
	}{
		{
			name:          "basic line",
			line:          "github.com/example/pkg/file.go:10.20,15.30 2 1",
			wantFileName:  "github.com/example/pkg/file.go",
			wantPackage:   "github.com/example/pkg",
			wantStatCount: 2,
			wantCallCount: 1,
		},
		{
			name:          "no coverage",
			line:          "github.com/example/pkg/file.go:10.20,15.30 5 0",
			wantFileName:  "github.com/example/pkg/file.go",
			wantPackage:   "github.com/example/pkg",
			wantStatCount: 5,
			wantCallCount: 0,
		},
		{
			name:          "deeply nested package",
			line:          "github.com/example/pkg/subpkg/subsubpkg/file.go:10.20,15.30 3 2",
			wantFileName:  "github.com/example/pkg/subpkg/subsubpkg/file.go",
			wantPackage:   "github.com/example/pkg/subpkg/subsubpkg",
			wantStatCount: 3,
			wantCallCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := new(Block)
			err := block.HydrateFromRawLine(tt.line)
			if err != nil {
				t.Errorf("HydrateFromRawLine() unexpected error = %v", err)
				return
			}

			if block.FileName != tt.wantFileName {
				t.Errorf("HydrateFromRawLine() FileName = %v, want %v",
					block.FileName, tt.wantFileName)
			}

			if block.PackageName != tt.wantPackage {
				t.Errorf("HydrateFromRawLine() PackageName = %v, want %v",
					block.PackageName, tt.wantPackage)
			}

			if block.StatementCount != tt.wantStatCount {
				t.Errorf("HydrateFromRawLine() StatementCount = %v, want %v",
					block.StatementCount, tt.wantStatCount)
			}

			if block.CallCount != tt.wantCallCount {
				t.Errorf("HydrateFromRawLine() CallCount = %v, want %v",
					block.CallCount, tt.wantCallCount)
			}

			if block.RawLineItem != tt.line {
				t.Errorf("HydrateFromRawLine() RawLineItem = %v, want %v",
					block.RawLineItem, tt.line)
			}

			// The BlockComponent should be the first part of the line before the statement and call counts
			parts := splitLast(tt.line, ' ', 2)
			expectedComponent := parts[0]
			if block.BlockComponent != expectedComponent {
				t.Errorf("HydrateFromRawLine() BlockComponent = %v, want %v",
					block.BlockComponent, expectedComponent)
			}
		})
	}
}

// Helper function to split a string by the last n occurrences of a separator
func splitLast(s string, sep byte, n int) []string {
	result := make([]string, 0, n+1)

	// Start from the end of the string
	remaining := s
	for range n {
		lastIdx := -1
		for j := len(remaining) - 1; j >= 0; j-- {
			if remaining[j] == sep {
				lastIdx = j
				break
			}
		}

		if lastIdx == -1 {
			// No more separators found
			break
		}

		// Add the part after the separator to the result
		result = append([]string{remaining[lastIdx+1:]}, result...)
		// Continue with the part before the separator
		remaining = remaining[:lastIdx]
	}

	// Add the remaining part
	result = append([]string{remaining}, result...)

	return result
}
