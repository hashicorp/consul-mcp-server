// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHumanVersion(t *testing.T) {
	tests := []struct {
		name              string
		version           string
		versionPrerelease string
		versionMetadata   string
		expectedContains  []string
	}{
		{
			name:              "release version",
			version:           "1.0.0",
			versionPrerelease: "",
			versionMetadata:   "",
			expectedContains:  []string{"1.0.0"},
		},
		{
			name:              "prerelease version",
			version:           "1.0.0",
			versionPrerelease: "beta1",
			versionMetadata:   "",
			expectedContains:  []string{"1.0.0", "beta1"},
		},
		{
			name:              "version with metadata",
			version:           "1.0.0",
			versionPrerelease: "",
			versionMetadata:   "20230901",
			expectedContains:  []string{"1.0.0", "20230901"},
		},
		{
			name:              "full version string",
			version:           "1.0.0",
			versionPrerelease: "rc1",
			versionMetadata:   "git.abc123",
			expectedContains:  []string{"1.0.0", "rc1", "git.abc123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily set version variables
			originalVersion := Version
			originalPrerelease := VersionPrerelease
			originalMetadata := VersionMetadata

			Version = tt.version
			VersionPrerelease = tt.versionPrerelease
			VersionMetadata = tt.versionMetadata

			defer func() {
				Version = originalVersion
				VersionPrerelease = originalPrerelease
				VersionMetadata = originalMetadata
			}()

			result := GetHumanVersion()

			// Check that all expected components are present
			for _, expected := range tt.expectedContains {
				assert.Contains(t, result, expected)
			}

			// Ensure version is not empty
			assert.NotEmpty(t, result)
		})
	}
}

func TestVersionConstants(t *testing.T) {
	t.Run("version constants are defined", func(t *testing.T) {
		// Test that version constants exist and have reasonable values
		assert.NotEmpty(t, Version, "Version should not be empty")

		// GitCommit and BuildDate might be empty in test environment
		// but should be strings
		assert.IsType(t, "", GitCommit)
		assert.IsType(t, "", BuildDate)
		assert.IsType(t, "", VersionPrerelease)
		assert.IsType(t, "", VersionMetadata)
	})
}

func TestVersionFormatting(t *testing.T) {
	t.Run("dev version formatting", func(t *testing.T) {
		originalVersion := Version
		originalPrerelease := VersionPrerelease

		Version = "dev"
		VersionPrerelease = ""

		defer func() {
			Version = originalVersion
			VersionPrerelease = originalPrerelease
		}()

		result := GetHumanVersion()
		assert.Contains(t, result, "dev")
	})

	t.Run("semantic version formatting", func(t *testing.T) {
		originalVersion := Version

		Version = "1.2.3"

		defer func() {
			Version = originalVersion
		}()

		result := GetHumanVersion()
		assert.Contains(t, result, "1.2.3")

		// Should follow semantic versioning pattern
		assert.Regexp(t, `\d+\.\d+\.\d+`, result)
	})
}

func TestVersionEdgeCases(t *testing.T) {
	t.Run("empty version", func(t *testing.T) {
		originalVersion := Version

		Version = ""

		defer func() {
			Version = originalVersion
		}()

		result := GetHumanVersion()
		// Should handle empty version gracefully
		assert.NotPanics(t, func() { GetHumanVersion() })
		assert.IsType(t, "", result)
	})

	t.Run("version with special characters", func(t *testing.T) {
		originalVersion := Version
		originalPrerelease := VersionPrerelease

		Version = "1.0.0"
		VersionPrerelease = "alpha.1+build.1"

		defer func() {
			Version = originalVersion
			VersionPrerelease = originalPrerelease
		}()

		result := GetHumanVersion()
		assert.Contains(t, result, "1.0.0")
		assert.Contains(t, result, "alpha.1+build.1")
	})
}
