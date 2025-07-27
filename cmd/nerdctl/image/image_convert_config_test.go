/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package image

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/containerd/nerdctl/v2/cmd/nerdctl/helpers"
	"github.com/containerd/nerdctl/v2/pkg/config"
)

func TestConvertOptionsWithConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir, err := os.MkdirTemp("", "nerdctl-config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "nerdctl.toml")
	configContent := `
[compression]
zstd_implementation = "klauspost"
zstd_compression_level = 11
zstd_chunked_compression_level = 22
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set up environment to use the test config
	oldEnv := os.Getenv("NERDCTL_TOML")
	os.Setenv("NERDCTL_TOML", configPath)
	defer os.Setenv("NERDCTL_TOML", oldEnv)

	// Create a test command
	cmd := &cobra.Command{}
	// Set up flags
	cmd.Flags().Bool("zstd", true, "")
	cmd.Flags().Int("zstd-compression-level", 3, "")
	cmd.Flags().Bool("zstdchunked", true, "")
	cmd.Flags().Int("zstdchunked-compression-level", 3, "")
	
	// Add other required flags
	cmd.Flags().String("format", "", "")
	cmd.Flags().Bool("estargz", false, "")
	cmd.Flags().String("estargz-record-in", "", "")
	cmd.Flags().Int("estargz-compression-level", 9, "")
	cmd.Flags().Int("estargz-chunk-size", 0, "")
	cmd.Flags().Int("estargz-min-chunk-size", 0, "")
	cmd.Flags().Bool("estargz-external-toc", false, "")
	cmd.Flags().Bool("estargz-keep-diff-id", false, "")
	cmd.Flags().String("zstdchunked-record-in", "", "")
	cmd.Flags().Int("zstdchunked-chunk-size", 0, "")
	cmd.Flags().Bool("nydus", false, "")
	cmd.Flags().String("nydus-builder-path", "", "")
	cmd.Flags().String("nydus-work-dir", "", "")
	cmd.Flags().String("nydus-prefetch-patterns", "", "")
	cmd.Flags().String("nydus-compressor", "", "")
	cmd.Flags().Bool("overlaybd", false, "")
	cmd.Flags().String("overlaybd-fs-type", "", "")
	cmd.Flags().String("overlaybd-dbstr", "", "")
	cmd.Flags().Bool("soci", false, "")
	cmd.Flags().Int64("soci-min-layer-size", -1, "")
	cmd.Flags().Int64("soci-span-size", -1, "")
	cmd.Flags().Bool("uncompress", false, "")
	cmd.Flags().Bool("oci", false, "")
	cmd.Flags().StringSlice("platform", []string{}, "")
	cmd.Flags().Bool("all-platforms", false, "")
	cmd.Flags().Bool("debug-compression", false, "")

	// Add global flags that ProcessRootCmdFlags expects
	cmd.PersistentFlags().Bool("debug", false, "")
	cmd.PersistentFlags().Bool("debug-full", false, "")
	cmd.PersistentFlags().String("address", "", "")
	cmd.PersistentFlags().String("namespace", "", "")
	cmd.PersistentFlags().String("snapshotter", "", "")
	cmd.PersistentFlags().String("cni-path", "", "")
	cmd.PersistentFlags().String("cni-netconfpath", "", "")
	cmd.PersistentFlags().String("data-root", "", "")
	cmd.PersistentFlags().String("cgroup-manager", "", "")
	cmd.PersistentFlags().Bool("insecure-registry", false, "")
	cmd.PersistentFlags().StringSlice("hosts-dir", []string{}, "")
	cmd.PersistentFlags().Bool("experimental", false, "")
	cmd.PersistentFlags().String("host-gateway-ip", "", "")
	cmd.PersistentFlags().String("bridge-ip", "", "")
	cmd.PersistentFlags().Bool("kube-hide-dupe", false, "")
	cmd.PersistentFlags().StringSlice("cdi-spec-dirs", []string{}, "")
	cmd.PersistentFlags().String("userns-remap", "", "")
	cmd.PersistentFlags().StringArray("global-dns", []string{}, "")
	cmd.PersistentFlags().StringArray("global-dns-opts", []string{}, "")
	cmd.PersistentFlags().StringArray("global-dns-search", []string{}, "")

	t.Run("ConfigDefaultsAreUsed", func(t *testing.T) {
		// Get options without setting CLI flags
		opts, err := convertOptions(cmd)
		require.NoError(t, err)

		// Config values should be used as defaults since flags weren't changed
		assert.Equal(t, 11, opts.ZstdCompressionLevel, "zstd compression level should use config default")
		assert.Equal(t, 22, opts.ZstdChunkedCompressionLevel, "zstdchunked compression level should use config default")
	})

	t.Run("CLIFlagsOverrideConfig", func(t *testing.T) {
		// Mark flags as changed to simulate CLI input
		cmd.Flags().Set("zstd-compression-level", "5")
		cmd.Flags().Set("zstdchunked-compression-level", "7")

		opts, err := convertOptions(cmd)
		require.NoError(t, err)

		// CLI values should override config
		assert.Equal(t, 5, opts.ZstdCompressionLevel, "CLI flag should override config")
		assert.Equal(t, 7, opts.ZstdChunkedCompressionLevel, "CLI flag should override config")
	})
}

func TestDebugCompressionFlag(t *testing.T) {
	// Create a test command
	cmd := convertCommand()
	
	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	
	// Test with debug flag
	cmd.SetArgs([]string{"--debug-compression", "--zstdchunked", "--oci", "test:src", "test:dst"})
	
	// We can't actually run the command because it requires a real containerd client
	// But we can verify the flag is registered
	debugFlag := cmd.Flag("debug-compression")
	assert.NotNil(t, debugFlag, "debug-compression flag should be registered")
	assert.Equal(t, "false", debugFlag.DefValue, "debug-compression default should be false")
}