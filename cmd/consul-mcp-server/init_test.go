// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitLoggerFilePermissions(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "mcp.log")

	logger, err := initLogger(logPath)
	require.NoError(t, err)
	t.Cleanup(func() { logger.Out.(*os.File).Close() })

	info, err := os.Stat(logPath)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
