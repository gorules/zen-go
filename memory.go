//go:build memory_test
// +build memory_test

package zen

// #cgo CPPFLAGS: -fsanitize=address
// #cgo LDFLAGS: -fsanitize=address
//
// #include <sanitizer/lsan_interface.h>
import "C"
import (
	"runtime"
)

// Call LLVM Leak Sanitizer's at-exit hook that doesn't
// get called automatically by Go.
func doLeakSanitizerCheck() {
	runtime.GC()
	C.__lsan_do_leak_check()
}
