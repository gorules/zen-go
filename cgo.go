package zen

/*
#cgo LDFLAGS: -pthread -lzen_ffi
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/deps/darwin_x86_64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/deps/darwin_arm64
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/deps/linux_x86_64 -ldl
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/deps/linux_arm64 -ldl
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/deps/linux_arm64 -ldl
*/
import "C"

import (
	_ "github.com/gorules/zen-go/deps/darwin_amd64"
	_ "github.com/gorules/zen-go/deps/darwin_arm64"
	_ "github.com/gorules/zen-go/deps/linux_amd64"
	_ "github.com/gorules/zen-go/deps/linux_arm64"
	_ "github.com/gorules/zen-go/deps/windows_amd64"
)
