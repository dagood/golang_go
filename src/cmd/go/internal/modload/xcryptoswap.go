// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modload

import (
	"cmd/go/internal/base"
	"cmd/go/internal/cfg"
	"cmd/go/internal/lockedfile"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

var (
	// True if x/crypto is being swapped out with a module replacement.
	xCryptoSwap bool
	// The path to the original go.mod file before x/crypto was swapped out.
	// Empty string if no replacement happened.
	preXCryptoSwapModFile string
)

// initXCryptoSwap determines whether x/crypto should be swapped out and if so, performs the swap by
// creating a temp go.mod with a new replace directive and passing it to cfg.ModFile.
func initXCryptoSwap() error {
	// Use build constraint evaluation to find if swapping is enabled in this build context.
	// This runs the build tag through the same system that is used to evaluate build tags in source
	// code, so we know it will match up with the build tags the dev intends to use.
	//
	// This looks weird: there is no "easy" API to check a build tag. We follow an approach also
	// used in go/build/build_test.go TestMatchFile: copy the build context to avoid modifying the
	// original and assign OpenFile to a new func that defines the content of a source file without
	// actually needing to write anything to disk.
	backendCheckContext := cfg.BuildContext
	backendCheckContext.OpenFile = func(path string) (io.ReadCloser, error) {
		source := "//go:build goexperiment.xcryptobackendswap"
		return io.NopCloser(strings.NewReader(source)), nil
	}
	match, err := backendCheckContext.MatchFile("", "backendcheck.go")
	if err != nil {
		return fmt.Errorf("error checking for crypto backend: %v", err)
	}
	if !match {
		return nil
	}

	switch cfg.CmdName {
	// Look for commands where we must not include the replace directive.
	//
	// As of writing, this follows the approach of setDefaultBuildMod and checks for specific
	// command names. When the TODO in setDefaultBuildMod is addressed to handle this more
	// elegantly (https://github.com/golang/go/issues/40775), we should follow suit.
	//
	// E.g. if we included the replacement with "go mod download", it could update the go.sum
	// file but without the checksum necessary to build with ordinary golang.org/x/crypto.
	case "get", "mod download", "mod init", "mod tidy", "mod edit", "mod vendor":
		return nil
	}
	if workFilePath != "" {
		return fmt.Errorf("unable to replace in workspace mode because '-modfile' is not supported")
	}
	// At this point there are no known reasons there would be 0 or >1 modRoots, but check just in
	// case we're missing a scenario.
	if len(modRoots) != 1 {
		return fmt.Errorf("expected 1 modroot, found %v, %#v", len(modRoots), modRoots)
	}

	oldFilePath := modFilePath(modRoots[0])
	data, err := lockedfile.Read(oldFilePath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod file: %v", err)
	}
	modFile, err := modfile.Parse(oldFilePath, data, nil)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod file: %v", err)
	}
	// If this is "go build -mod=mod", include the replacement even if the go.mod doesn't include
	// x/crypto. There might be a transitive dependency we don't know about until imports are
	// evaluated and find an x/crypto reference.
	if !modFileUsesXCrypto(modFile) && cfg.BuildMod != "mod" {
		return nil
	}
	// Add (or update!) a replace directive.
	modFile.AddReplace(
		"golang.org/x/crypto",
		"",
		xCryptoNewModPath(),
		"",
	)
	modFile.Cleanup()
	afterData, err := modFile.Format()
	if err != nil {
		return fmt.Errorf("failed to format (marshal) go.mod file: %v", err)
	}
	f, err := os.CreateTemp("", "xcrypto-replacement-go-*.mod")
	if err != nil {
		return fmt.Errorf("failed to create temp go.mod file: %v", err)
	}
	base.AtExit(func() {
		// Best effort cleanup.
		if err := os.Remove(f.Name()); err != nil {
			base.Errorf("go: failed to remove temp go.mod file %q with x/crypto replacement: %v", f.Name(), err)
		}
	})
	if _, err := f.Write(afterData); err != nil {
		return fmt.Errorf("failed to write temp go.mod file data: %v", err)
	}
	cfg.ModFile = f.Name()
	xCryptoSwap = true
	preXCryptoSwapModFile = oldFilePath
	return nil
}

func modFileUsesXCrypto(m *modfile.File) bool {
	for _, r := range m.Require {
		if r.Mod.Path == "golang.org/x/crypto" {
			return true
		}
	}
	return false
}

// toolsetOverrideModuleDir returns the path to a directory in GOROOT where copies of modules are
// kept that may be used when building a program with that particular toolset. It is similar to a
// vendor directory, but only contains the modules, no modules.txt.
func toolsetOverrideModuleDir() string {
	return filepath.Join(cfg.BuildContext.GOROOT, "ms_mod")
}

// xCryptoNewModPath returns the module path that x/crypto should be replaced with.
func xCryptoNewModPath() string {
	return filepath.Join(toolsetOverrideModuleDir(), "golang.org", "x", "crypto")
}

func toolsetOverridePackage(path string) (dir string, vendorDir string, ok bool) {
	if !HasModRoot() || !xCryptoSwap {
		return "", "", false
	}
	vendorDir = toolsetOverrideModuleDir()
	// Before checking the main vendor directory, check the toolset's override vendor directory.
	// This is where the x/crypto fork is located that needs to be swapped in.
	dir, haveGoFiles, _ := dirInModule(path, "", vendorDir, false)
	if !haveGoFiles {
		return "", "", false
	}
	return dir, vendorDir, true
}
