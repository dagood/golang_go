// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mstoolset implements Microsoft Go toolset customizations that let a
// fork of the standard library be built against an out-of-tree vendor
// directory, go.mod, and go.sum without re-vendoring GOROOT/src on every
// change.
//
// The following environment variables are recognized:
//
//	MS_GOTOOLSETVENDOR  Use this directory in place of GOROOT/src/vendor.
//	MS_GOTOOLSETGOMOD   Use this file in place of GOROOT/src/go.mod.
//	MS_GOTOOLSETGOSUM   Use this file in place of GOROOT/src/go.sum.
//
// The redirection is implemented through the cmd/go virtual file system
// (cmd/go/internal/fsys), the same mechanism used by the FIPS 140 snapshot, so
// every part of the go command that reads those paths transparently sees the
// override. The variables only affect the "std" module rooted at GOROOT/src;
// they are designed to work during the bootstrap performed by make.bash.
package mstoolset

import (
	"os"
	"path/filepath"

	"cmd/go/internal/cfg"
	"cmd/go/internal/fsys"
)

// Environment variable names.
const (
	VendorEnv = "MS_GOTOOLSETVENDOR"
	GoModEnv  = "MS_GOTOOLSETGOMOD"
	GoSumEnv  = "MS_GOTOOLSETGOSUM"
)

var initDone bool

// Init registers the MS_GOTOOLSET* overrides with the virtual file system.
//
// It must be called after the working directory is final (so that relative
// override paths resolve correctly) and before fsys.Init, so the registered
// replacements are incorporated into the overlay. Init is idempotent.
func Init() {
	if initDone {
		return
	}
	initDone = true

	// cfg.GOROOTsrc is GOROOT/src, the root of the "std" module.
	srcDir := cfg.GOROOTsrc
	if srcDir == "" {
		return
	}

	if dir := os.Getenv(VendorEnv); dir != "" {
		fsys.Bind(dir, filepath.Join(srcDir, "vendor"))
	}
	if file := os.Getenv(GoModEnv); file != "" {
		fsys.Replace(filepath.Join(srcDir, "go.mod"), file)
	}
	if file := os.Getenv(GoSumEnv); file != "" {
		fsys.Replace(filepath.Join(srcDir, "go.sum"), file)
	}
}
