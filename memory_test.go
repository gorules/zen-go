//go:build memory_test
// +build memory_test

package zen

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	doLeakSanitizerCheck()
	os.Exit(exitCode)
}
