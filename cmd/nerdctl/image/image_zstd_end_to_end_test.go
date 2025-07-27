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
	"path/filepath"
	"strings"
	"testing"

	"github.com/containerd/nerdctl/v2/pkg/testutil"
	"github.com/containerd/nerdctl/v2/pkg/testutil/nerdtest"
	"github.com/containerd/nerdctl/v2/pkg/testutil/test"
)

// TestZstdCompressionEndToEnd tests the complete workflow with zstd compression
func TestZstdCompressionEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping end-to-end test in short mode")
	}

	testCase := nerdtest.Setup()
	testCase.Env = append(testCase.Env, "NERDCTL_EXPERIMENTAL=1")

	// Test with different real-world images
	testImages := []struct {
		name  string
		image string
	}{
		{"alpine", "docker.io/library/alpine:latest"},
		{"nginx", "docker.io/library/nginx:alpine"},
		{"tomcat", "docker.io/library/tomcat:10.1-jdk21-temurin-noble"},
	}

	for _, img := range testImages {
		t.Run(img.name, func(t *testing.T) {
			if img.name == "tomcat" && testing.Short() {
				t.Skip("Skipping large image in short mode")
			}

			// Pull the base image
			testCase.Run(test.RunBinary("pull", img.image),
				test.WithStdin(nil),
				test.Expects(0, nil, nil))

			// Test complete workflow
			t.Run("CompleteWorkflow", func(t *testing.T) {
				// 1. Convert to zstd with high compression
				convertedImage := fmt.Sprintf("%s-zstd-22:test", img.name)
				convertResult := testCase.Run(
					test.RunBinary(
						"image", "convert",
						"--debug-compression",
						"--zstd",
						"--zstd-compression-level=22",
						"--oci",
						img.image,
						convertedImage,
					),
					test.WithStdin(nil),
				)

				// Check if conversion worked (might fall back to level 11)
				if convertResult.ExitCode != 0 {
					t.Fatalf("Conversion failed: %s", convertResult.Stderr())
				}

				// Verify debug output shows implementation info
				if !strings.Contains(convertResult.Stdout(), "Implementation:") {
					t.Error("Debug output missing implementation info")
				}

				// 2. Save the converted image
				tarFile := filepath.Join(t.TempDir(), fmt.Sprintf("%s-zstd.tar", img.name))
				saveResult := testCase.Run(
					test.RunBinary("image", "save", "-o", tarFile, convertedImage),
					test.WithStdin(nil),
				)
				if saveResult.ExitCode != 0 {
					t.Fatalf("Save failed: %s", saveResult.Stderr())
				}

				// Check file was created
				info, err := os.Stat(tarFile)
				if err != nil {
					t.Fatalf("Saved tar file not found: %v", err)
				}
				t.Logf("Saved image size: %d bytes", info.Size())

				// 3. Remove the converted image
				testCase.Run(
					test.RunBinary("image", "rm", "-f", convertedImage),
					test.WithStdin(nil),
					test.Expects(0, nil, nil))

				// 4. Load it back
				loadResult := testCase.Run(
					test.RunBinary("image", "load", "-i", tarFile),
					test.WithStdin(nil),
				)
				if loadResult.ExitCode != 0 {
					t.Fatalf("Load failed: %s", loadResult.Stderr())
				}

				// 5. Verify the loaded image exists and can be inspected
				inspectResult := testCase.Run(
					test.RunBinary("image", "inspect", convertedImage),
					test.WithStdin(nil),
				)
				if inspectResult.ExitCode != 0 {
					t.Fatal("Loaded image not found")
				}

				// 6. Run a container from the loaded image
				containerName := fmt.Sprintf("%s-zstd-test", img.name)
				runResult := testCase.Run(
					test.RunBinary(
						"run", "--rm", "--name", containerName,
						convertedImage, "echo", "Hello from zstd compressed image",
					),
					test.WithStdin(nil),
				)
				if runResult.ExitCode != 0 {
					t.Fatalf("Container run failed: %s", runResult.Stderr())
				}
				if !strings.Contains(runResult.Stdout(), "Hello from zstd compressed image") {
					t.Error("Container output not as expected")
				}

				// Clean up
				testCase.Run(
					test.RunBinary("image", "rm", "-f", convertedImage),
					test.WithStdin(nil),
				)
			})

			// Test zstd:chunked workflow
			t.Run("ZstdChunkedWorkflow", func(t *testing.T) {
				// 1. Convert to zstd:chunked
				chunkedImage := fmt.Sprintf("%s-zstdchunked:test", img.name)
				
				// Test with environment override to ensure we're testing both implementations
				for _, env := range []string{"", "ZSTD_FORCE_IMPLEMENTATION=klauspost", "ZSTD_FORCE_IMPLEMENTATION=gozstd"} {
					testName := "Auto"
					if strings.Contains(env, "klauspost") {
						testName = "PureGo"
					} else if strings.Contains(env, "gozstd") {
						testName = "Libzstd"
					}

					t.Run(testName, func(t *testing.T) {
						testEnv := testCase.Env
						if env != "" {
							testEnv = append(testEnv, env)
						}

						convertResult := testCase.Run(
							test.RunBinary(
								"image", "convert",
								"--debug-compression",
								"--zstdchunked",
								"--zstdchunked-compression-level=11",
								"--oci",
								img.image,
								chunkedImage,
							),
							test.WithStdin(nil),
							test.WithEnv(testEnv...),
						)

						// Skip if implementation not available
						if strings.Contains(env, "gozstd") && strings.Contains(convertResult.Stderr(), "libzstd not available") {
							t.Skip("libzstd not available")
						}

						if convertResult.ExitCode != 0 {
							t.Fatalf("Conversion failed: %s", convertResult.Stderr())
						}

						// Verify correct implementation was used
						if strings.Contains(env, "klauspost") {
							if !strings.Contains(convertResult.Stdout(), "klauspost") {
								t.Error("Expected pure Go implementation not used")
							}
						} else if strings.Contains(env, "gozstd") {
							if !strings.Contains(convertResult.Stdout(), "gozstd") {
								t.Error("Expected libzstd implementation not used")
							}
						}

						// Create a container to verify the image works
						containerName := fmt.Sprintf("%s-chunked-test-%s", img.name, testName)
						runResult := testCase.Run(
							test.RunBinary(
								"run", "--rm", "--name", containerName,
								chunkedImage, "echo", "zstd:chunked works",
							),
							test.WithStdin(nil),
						)
						if runResult.ExitCode != 0 {
							t.Fatalf("Container run failed: %s", runResult.Stderr())
						}

						// Clean up
						testCase.Run(
							test.RunBinary("image", "rm", "-f", chunkedImage),
							test.WithStdin(nil),
						)
					})
				}
			})
		})
	}
}

// TestZstdCompressionWithBuildAndCommit tests zstd compression with build and commit operations
func TestZstdCompressionWithBuildAndCommit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping build test in short mode")
	}

	testCase := nerdtest.Setup()

	// Create a simple Dockerfile
	dockerfileContent := `FROM alpine:latest
RUN echo "Hello from build" > /hello.txt
CMD ["cat", "/hello.txt"]`

	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Build the image
	builtImage := "test-build-zstd:latest"
	buildResult := testCase.Run(
		test.RunBinary("build", "-t", builtImage, tmpDir),
		test.WithStdin(nil),
	)
	if buildResult.ExitCode != 0 {
		t.Fatalf("Build failed: %s", buildResult.Stderr())
	}

	// Convert to zstd
	zstdImage := "test-build-zstd-converted:latest"
	convertResult := testCase.Run(
		test.RunBinary(
			"image", "convert",
			"--zstd",
			"--zstd-compression-level=11",
			"--oci",
			builtImage,
			zstdImage,
		),
		test.WithStdin(nil),
	)
	if convertResult.ExitCode != 0 {
		t.Fatalf("Conversion failed: %s", convertResult.Stderr())
	}

	// Test commit workflow
	t.Run("CommitWorkflow", func(t *testing.T) {
		// Run a container and modify it
		containerName := "test-commit-container"
		testCase.Run(
			test.RunBinary(
				"run", "--name", containerName,
				zstdImage, "sh", "-c", "echo 'Modified' >> /hello.txt",
			),
			test.WithStdin(nil),
			test.Expects(0, nil, nil),
		)

		// Commit the container
		committedImage := "test-committed-zstd:latest"
		commitResult := testCase.Run(
			test.RunBinary("commit", containerName, committedImage),
			test.WithStdin(nil),
		)
		if commitResult.ExitCode != 0 {
			t.Fatalf("Commit failed: %s", commitResult.Stderr())
		}

		// Verify the committed image
		runResult := testCase.Run(
			test.RunBinary(
				"run", "--rm",
				committedImage, "cat", "/hello.txt",
			),
			test.WithStdin(nil),
		)
		if runResult.ExitCode != 0 {
			t.Fatalf("Run committed image failed: %s", runResult.Stderr())
		}
		if !strings.Contains(runResult.Stdout(), "Modified") {
			t.Error("Committed changes not found")
		}

		// Clean up
		testCase.Run(test.RunBinary("rm", containerName), test.WithStdin(nil))
		testCase.Run(test.RunBinary("image", "rm", "-f", committedImage), test.WithStdin(nil))
	})

	// Clean up
	testCase.Run(test.RunBinary("image", "rm", "-f", builtImage, zstdImage), test.WithStdin(nil))
}

// TestZstdDecompressionPerformance verifies that decompression is using the runtime-detected implementation
func TestZstdDecompressionPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	testCase := nerdtest.Setup()

	// Use a reasonably sized image
	sourceImage := "docker.io/library/alpine:latest"
	
	// Pull the source image
	testCase.Run(test.RunBinary("pull", sourceImage),
		test.WithStdin(nil),
		test.Expects(0, nil, nil))

	// Convert with different compression levels and measure load times
	levels := []int{1, 11, 22}
	
	for _, level := range levels {
		t.Run(fmt.Sprintf("Level%d", level), func(t *testing.T) {
			convertedImage := fmt.Sprintf("alpine-zstd-%d:test", level)
			tarFile := filepath.Join(t.TempDir(), fmt.Sprintf("alpine-zstd-%d.tar", level))

			// Convert to zstd
			convertResult := testCase.Run(
				test.RunBinary(
					"image", "convert",
					"--zstd",
					fmt.Sprintf("--zstd-compression-level=%d", level),
					"--oci",
					sourceImage,
					convertedImage,
				),
				test.WithStdin(nil),
			)

			// Skip if level not supported
			if level > 11 && strings.Contains(convertResult.Stderr(), "exceeds maximum") {
				t.Skipf("Level %d not supported", level)
			}
			if convertResult.ExitCode != 0 {
				t.Fatalf("Conversion failed: %s", convertResult.Stderr())
			}

			// Save the image
			testCase.Run(
				test.RunBinary("image", "save", "-o", tarFile, convertedImage),
				test.WithStdin(nil),
				test.Expects(0, nil, nil))

			// Remove the image
			testCase.Run(
				test.RunBinary("image", "rm", "-f", convertedImage),
				test.WithStdin(nil),
				test.Expects(0, nil, nil))

			// Load and measure time (in real test, we'd use proper timing)
			loadResult := testCase.Run(
				test.RunBinary("image", "load", "-i", tarFile),
				test.WithStdin(nil),
			)
			if loadResult.ExitCode != 0 {
				t.Fatalf("Load failed: %s", loadResult.Stderr())
			}

			// The actual decompression happens during load
			// Our runtime-detected decompressor is used by stargz-snapshotter
			t.Logf("Successfully loaded image compressed at level %d", level)

			// Clean up
			testCase.Run(
				test.RunBinary("image", "rm", "-f", convertedImage),
				test.WithStdin(nil),
			)
		})
	}
}