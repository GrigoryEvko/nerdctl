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

package zstd

import (
	"os"
	"testing"
	
	"github.com/containerd/log"
	"github.com/sirupsen/logrus"
)

// TestMain sets up the test environment for all tests in this package
func TestMain(m *testing.M) {
	// Set single-threaded mode for all tests to ensure deterministic behavior
	if os.Getenv("ZSTD_WORKERS") == "" {
		os.Setenv("ZSTD_WORKERS", "1")
	}

	// Initialize the logger to prevent nil pointer issues
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Only show errors in tests
	log.L = logrus.NewEntry(logger)

	// Run tests
	code := m.Run()
	
	os.Exit(code)
}