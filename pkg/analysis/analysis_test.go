package analysis

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name               string
		coverProfileData   string
		coverageStandard   float32
		wantErr            bool
		wantCoverage       float32
		wantStandardMet    bool
		wantPackageCount   int
		wantBlockCount     int
		wantBlockCallCount int
	}{
		{
			name: "valid profile with full coverage",
			coverProfileData: `mode: set
github.com/example/pkg/file1.go:10.20,15.30 2 1
github.com/example/pkg/file1.go:15.30,20.40 2 1
github.com/example/pkg/file1.go:20.40,25.50 2 1`,
			coverageStandard:   80.0,
			wantErr:            false,
			wantCoverage:       100.0,
			wantStandardMet:    true,
			wantPackageCount:   1,
			wantBlockCount:     3,
			wantBlockCallCount: 3,
		},
		{
			name: "valid profile with partial coverage",
			coverProfileData: `mode: set
github.com/example/pkg/file1.go:10.20,15.30 2 1
github.com/example/pkg/file1.go:15.30,20.40 2 0
github.com/example/pkg/file1.go:20.40,25.50 2 1`,
			coverageStandard:   80.0,
			wantErr:            false,
			wantCoverage:       66.66666,
			wantStandardMet:    false,
			wantPackageCount:   1,
			wantBlockCount:     3,
			wantBlockCallCount: 2,
		},
		{
			name: "multiple packages with different coverage",
			coverProfileData: `mode: set
github.com/example/pkg1/file1.go:10.20,15.30 2 1
github.com/example/pkg1/file1.go:15.30,20.40 2 1
github.com/example/pkg2/file2.go:10.20,15.30 2 0
github.com/example/pkg2/file2.go:15.30,20.40 2 0`,
			coverageStandard:   50.0,
			wantErr:            false,
			wantCoverage:       50.0,
			wantStandardMet:    true,
			wantPackageCount:   2,
			wantBlockCount:     4,
			wantBlockCallCount: 2,
		},
		{
			name:               "empty profile",
			coverProfileData:   `mode: set`,
			coverageStandard:   80.0,
			wantErr:            false,
			wantCoverage:       0.0,
			wantStandardMet:    false,
			wantPackageCount:   0,
			wantBlockCount:     0,
			wantBlockCallCount: 0,
		},
		{
			name: "duplicate blocks with accumulated call counts",
			coverProfileData: `mode: set
github.com/example/pkg/file1.go:10.20,15.30 2 1
github.com/example/pkg/file1.go:10.20,15.30 2 2`,
			coverageStandard:   80.0,
			wantErr:            false,
			wantCoverage:       100.0,
			wantStandardMet:    true,
			wantPackageCount:   1,
			wantBlockCount:     1,
			wantBlockCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the test data
			tempFile, err := os.CreateTemp(t.TempDir(), "coverage-*.out")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			//nolint: errcheck
			defer os.Remove(tempFile.Name())

			_, err = tempFile.WriteString(tt.coverProfileData)
			if err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			//nolint: errcheck,gosec
			tempFile.Close()

			// Run the analysis
			summary, err := Run(tempFile.Name(), tt.coverageStandard)

			// Check error status
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Validate summary fields
			if summary.PackageCount != tt.wantPackageCount {
				t.Errorf("Run() PackageCount = %v, want %v", summary.PackageCount, tt.wantPackageCount)
			}

			if summary.BlockCount != tt.wantBlockCount {
				t.Errorf("Run() BlockCount = %v, want %v", summary.BlockCount, tt.wantBlockCount)
			}

			if summary.BlockCallCount != tt.wantBlockCallCount {
				t.Errorf("Run() BlockCallCount = %v, want %v", summary.BlockCallCount, tt.wantBlockCallCount)
			}

			if summary.BlocksNotCoveredCount != (tt.wantBlockCount - tt.wantBlockCallCount) {
				t.Errorf("Run() BlocksNotCoveredCount = %v, want %v",
					summary.BlocksNotCoveredCount, (tt.wantBlockCount - tt.wantBlockCallCount))
			}

			// Check coverage percentage within small delta for floating point comparison
			delta := float32(0.01)
			if abs(summary.CoveragePercentage-tt.wantCoverage) > delta {
				t.Errorf("Run() CoveragePercentage = %v, want %v (±%v)",
					summary.CoveragePercentage, tt.wantCoverage, delta)
			}

			if summary.CoverageStandard != tt.coverageStandard {
				t.Errorf("Run() CoverageStandard = %v, want %v",
					summary.CoverageStandard, tt.coverageStandard)
			}

			if summary.CoverageStandardMet != tt.wantStandardMet {
				t.Errorf("Run() CoverageStandardMet = %v, want %v",
					summary.CoverageStandardMet, tt.wantStandardMet)
			}
		})
	}
}

func TestRun_FileNotFound(t *testing.T) {
	_, err := Run("nonexistent-file.out", 80.0)
	if err == nil {
		t.Error("Run() with nonexistent file should return error")
	}
}

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want float32
	}{
		{"normal division", 5, 2, 2.5},
		{"zero numerator", 0, 5, 0},
		{"zero denominator", 5, 0, 0},
		{"zero both", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := safeDivide(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("safeDivide() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to calculate absolute difference between float32 values
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
