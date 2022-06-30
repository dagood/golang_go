// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backend

import (
	"crypto/internal/boring/sig"
	"runtime"
	"syscall"
)

func init() {
	if v, r, ok := envGoFIPS(); ok && v == "1" {
		if !Enabled {
			if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
				panic("FIPS mode requested (" + r + ") but no crypto backend is supported on " + runtime.GOOS)
			}
			panic("FIPS mode requested (" + r + ") but no supported crypto backend is enabled")
		}
	}
}

func envGoFIPS() (value string, reason string, ok bool) {
	// TODO: Decide which environment variable to use.
	// See https://github.com/microsoft/go/issues/397.
	var varName string
	if value, ok = syscall.Getenv("GOFIPS"); ok {
		varName = "GOFIPS"
	} else if value, ok = syscall.Getenv("GOLANG_FIPS"); ok {
		varName = "GOLANG_FIPS"
	}
	if isRequireFIPS {
		if ok && value != "1" {
			panic("the 'requirefips' build tag is enabled, but it conflicts " +
				"with the detected env variable " +
				varName + "=" + value +
				" which would disable FIPS mode")
		}
		return "1", "requirefips tag set", true
	}
	if ok {
		return value, "environment variable " + varName + "=1", true
	}
	return "", "", false
}

// Unreachable marks code that should be unreachable
// when backend is in use.
func Unreachable() {
	if Enabled {
		panic("cryptobackend: invalid code execution")
	} else {
		// Code that's unreachable is exactly the code
		// we want to detect for reporting standard Go crypto.
		sig.StandardCrypto()
	}
}

// Provided by runtime.crypto_backend_runtime_arg0 to avoid os import.
func runtime_arg0() string

func hasSuffix(s, t string) bool {
	return len(s) > len(t) && s[len(s)-len(t):] == t
}

// UnreachableExceptTests marks code that should be unreachable
// when backend is in use. It panics.
func UnreachableExceptTests() {
	if Enabled {
		name := runtime_arg0()
		// If ran on Windows we'd need to allow _test.exe and .test.exe as well.
		if !hasSuffix(name, "_test") && !hasSuffix(name, ".test") {
			println("cryptobackend: unexpected code execution in", name)
			panic("cryptobackend: invalid code execution")
		}
	}
}
