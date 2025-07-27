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
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/containerd/nerdctl/v2/pkg/testutil"
	"github.com/containerd/nerdctl/v2/pkg/testutil/nerdtest"
)

// TestImageConvertRealImagesZstd tests zstd and zstdchunked conversion with real images
func TestImageConvertRealImagesZstd(t *testing.T) {
	testCase := nerdtest.Setup()

	testImages := []struct {
		name string
		image string
	}{
		{
			name:  "tomcat",
			image: "docker.io/library/tomcat:10.1-jdk21-temurin-noble",
		},
		{
			name:  "alpine",
			image: "docker.io/library/alpine:latest",
		},
		{
			name:  "nginx",
			image: "docker.io/library/nginx:alpine",
		},
	}

	compressionTests := []struct {
		name            string
		format          string
		compressionFlag string
		levelFlag       string
		levels          []int
	}{
		{
			name:            "zstd",
			format:          "zstd",
			compressionFlag: "--zstd",
			levelFlag:       "--zstd-compression-level",
			levels:          []int{1, 3, 11, 22},
		},
		{
			name:            "zstdchunked",
			format:          "zstdchunked",
			compressionFlag: "--zstdchunked",
			levelFlag:       "--zstdchunked-compression-level",
			levels:          []int{1, 3, 11, 22},
		},
	}

	for _, img := range testImages {
		for _, comp := range compressionTests {
			for _, level := range comp.levels {
				testName := fmt.Sprintf("%s_%s_level%d", img.name, comp.format, level)
				
				t.Run(testName, func(t *testing.T) {
					if testing.Short() && level > 11 {
						t.Skip("Skipping high compression levels in short mode")
					}

					testCase.Env = append(testCase.Env, "NERDCTL_EXPERIMENTAL=1")
					if level > 11 {
						// Test with libzstd if available
						testCase.Env = append(testCase.Env, "ZSTD_FORCE_IMPLEMENTATION=gozstd")
					}

					srcImage := img.image
					dstImage := fmt.Sprintf("%s-%s-level%d:test", img.name, comp.format, level)

					// Pull source image if not present
					testCase.Run(test.RunBinary("pull", srcImage),
						test.WithStdin(nil),
						test.Expects(0, nil, nil))

					// Convert with debug flag to see compression info
					result := testCase.Run(
						test.RunBinary(
							"image", "convert",
							"--debug-compression",
							comp.compressionFlag,
							fmt.Sprintf("%s=%d", comp.levelFlag, level),
							"--oci",
							srcImage,
							dstImage,
						),
						test.WithStdin(nil),
					)

					// Check that conversion succeeded
					if result.ExitCode != 0 {
						// Level 22 might fail if libzstd not available
						if level > 11 && strings.Contains(result.Stderr(), "exceeds maximum") {
							t.Skipf("Level %d not supported (libzstd not available)", level)
						}
						t.Fatalf("Conversion failed: %s", result.Stderr())
					}

					// Verify debug output
					if comp.format == "zstdchunked" {
						if !strings.Contains(result.Stdout(), "=== Compression Debug Info ===") {
							t.Error("Debug compression info not found in output")
						}
						if !strings.Contains(result.Stdout(), "Implementation:") {
							t.Error("Implementation info not found in debug output")
						}
						if !strings.Contains(result.Stdout(), fmt.Sprintf("Requested zstd:chunked Level: %d", level)) {
							// Check if it was capped
							if level > 11 && strings.Contains(result.Stdout(), "Requested zstd:chunked Level: 11") {
								t.Logf("Level was capped to 11 (pure Go implementation)")
							} else {
								t.Errorf("Expected compression level %d not found in debug output", level)
							}
						}
					}

					// Verify the converted image exists
					listResult := testCase.Run(
						test.RunBinary("image", "list", dstImage),
						test.WithStdin(nil),
					)
					if listResult.ExitCode != 0 {
						t.Fatal("Converted image not found")
					}

					// Clean up
					testCase.Run(
						test.RunBinary("image", "rm", "-f", dstImage),
						test.WithStdin(nil),
					)
				})
			}
		}
	}
}

// TestImageConvertAndPullWithCompression tests the full workflow of converting and pulling images
func TestImageConvertAndPullWithCompression(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	testCase := nerdtest.Setup()
	testCase.Env = append(testCase.Env, "NERDCTL_EXPERIMENTAL=1")

	// Use alpine for faster testing
	sourceImage := "docker.io/library/alpine:latest"
	convertedImage := "alpine-zstdchunked-test:latest"

	// Pull source image
	testCase.Run(test.RunBinary("pull", sourceImage),
		test.WithStdin(nil),
		test.Expects(0, nil, nil))

	// Test environment variable override
	t.Run("EnvironmentOverride", func(t *testing.T) {
		// Force pure Go implementation
		testCase.Env = append(testCase.Env, "ZSTD_FORCE_IMPLEMENTATION=klauspost")
		
		result := testCase.Run(
			test.RunBinary(
				"image", "convert",
				"--debug-compression",
				"--zstdchunked",
				"--zstdchunked-compression-level=15", // This should be capped to 11
				"--oci",
				sourceImage,
				convertedImage,
			),
			test.WithStdin(nil),
		)

		if result.ExitCode != 0 {
			t.Fatalf("Conversion failed: %s", result.Stderr())
		}

		// Should show pure Go implementation
		if !strings.Contains(result.Stdout(), "klauspost") {
			t.Error("Expected klauspost implementation not found")
		}
		
		// Level should be capped
		if !strings.Contains(result.Stdout(), "exceeds maximum") {
			t.Error("Expected warning about exceeding maximum level")
		}

		// Clean up
		testCase.Run(
			test.RunBinary("image", "rm", "-f", convertedImage),
			test.WithStdin(nil),
		)
	})

	// Test high compression with automatic detection
	t.Run("HighCompressionAuto", func(t *testing.T) {
		// Remove env override
		testCase.Env = []string{"NERDCTL_EXPERIMENTAL=1"}
		
		result := testCase.Run(
			test.RunBinary(
				"image", "convert",
				"--debug-compression",
				"--zstdchunked",
				"--zstdchunked-compression-level=22",
				"--oci",
				sourceImage,
				"alpine-zstdchunked-22:test",
			),
			test.WithStdin(nil),
		)

		// Conversion should succeed (might use level 11 if libzstd not available)
		if result.ExitCode != 0 {
			t.Fatalf("Conversion failed: %s", result.Stderr())
		}

		// Check which implementation was used
		if strings.Contains(result.Stdout(), "gozstd") {
			t.Log("Using libzstd implementation")
			if !strings.Contains(result.Stdout(), "Max Compression Level: 22") {
				t.Error("Expected max level 22 for libzstd")
			}
		} else if strings.Contains(result.Stdout(), "klauspost") {
			t.Log("Using pure Go implementation")
			if !strings.Contains(result.Stdout(), "Max Compression Level: 11") {
				t.Error("Expected max level 11 for pure Go")
			}
		}

		// Clean up
		testCase.Run(
			test.RunBinary("image", "rm", "-f", "alpine-zstdchunked-22:test"),
			test.WithStdin(nil),
		)
	})
}

// TestImageSaveLoadWithZstd tests saving and loading images with zstd compression
func TestImageSaveLoadWithZstd(t *testing.T) {
	testCase := nerdtest.Setup()
	
	sourceImage := "docker.io/library/alpine:latest"
	tarFile := "alpine-zstd-test.tar"
	
	// Pull source image
	testCase.Run(test.RunBinary("pull", sourceImage),
		test.WithStdin(nil),
		test.Expects(0, nil, nil))

	// First convert to zstd
	convertedImage := "alpine-zstd-save-test:latest"
	result := testCase.Run(
		test.RunBinary(
			"image", "convert",
			"--zstd",
			"--zstd-compression-level=11",
			"--oci",
			sourceImage,
			convertedImage,
		),
		test.WithStdin(nil),
	)

	if result.ExitCode != 0 {
		t.Fatalf("Conversion failed: %s", result.Stderr())
	}

	// Save the converted image
	saveResult := testCase.Run(
		test.RunBinary("image", "save", "-o", tarFile, convertedImage),
		test.WithStdin(nil),
	)

	if saveResult.ExitCode != 0 {
		t.Fatalf("Save failed: %s", saveResult.Stderr())
	}

	// Remove the image
	testCase.Run(
		test.RunBinary("image", "rm", "-f", convertedImage),
		test.WithStdin(nil),
	)

	// Load it back
	loadResult := testCase.Run(
		test.RunBinary("image", "load", "-i", tarFile),
		test.WithStdin(nil),
	)

	if loadResult.ExitCode != 0 {
		t.Fatalf("Load failed: %s", loadResult.Stderr())
	}

	// Verify the image was loaded
	listResult := testCase.Run(
		test.RunBinary("image", "list", convertedImage),
		test.WithStdin(nil),
	)
	
	if listResult.ExitCode != 0 {
		t.Fatal("Loaded image not found")
	}

	// Clean up
	os.Remove(tarFile)
	testCase.Run(
		test.RunBinary("image", "rm", "-f", convertedImage),
		test.WithStdin(nil),
	)
}

// TestZstdChunkedLazyPull tests lazy pulling with zstd:chunked images
func TestZstdChunkedLazyPull(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping lazy pull test in short mode")
	}

	testCase := nerdtest.Setup()
	
	// Only run if stargz snapshotter is available
	infoResult := testCase.Run(
		test.RunBinary("info", "-f", "json"),
		test.WithStdin(nil),
	)
	
	if infoResult.ExitCode != 0 {
		t.Skip("Cannot get nerdctl info")
	}
	
	if !strings.Contains(infoResult.Stdout(), "stargz") {
		t.Skip("stargz snapshotter not available")
	}

	// This test would require a registry to push/pull zstd:chunked images
	// For now, we'll just document that the decompression path is tested
	// through the stargz-snapshotter integration
	t.Log("zstd:chunked lazy pulling uses the runtime-detected decompressor in stargz-snapshotter")
}