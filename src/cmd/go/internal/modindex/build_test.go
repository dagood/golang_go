// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is not a copy. It tests that the copied code in this directory is
// maintained. Specifically, it tests areas that are modified by microsoft/go.
// The files also contain intentional modifications, so it isn't reasonable (as
// of writing) to test that the entire file is identical.

package modindex

import (
	"flag"
	"os"
	"strings"
	"testing"
)

var fixCopy = flag.Bool("fixcopy", false, "if true, update some copied code in build.go")

func TestCopyIdentical(t *testing.T) {
	originalBytes, err := os.ReadFile("../../../../go/build/build.go")
	if err != nil {
		t.Fatal(err)
	}
	wantCode := string(originalBytes)

	gotBytes, err := os.ReadFile("build.go")
	if err != nil {
		t.Fatal(err)
	}
	gotCode := string(gotBytes)

	tests := []struct {
		name   string
		prefix string
		suffix string
	}{
		{"matchTag", "func (ctxt *Context) matchTag(name string, allTags map[string]bool) bool {", "\n}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var want, got string
			if _, after, ok := strings.Cut(wantCode, tt.prefix); ok {
				if before, _, ok := strings.Cut(after, tt.suffix); ok {
					want = before
				} else {
					t.Fatal("suffix not found in original file")
				}
			} else {
				t.Fatal("prefix not found in original file")
			}
			if _, after, ok := strings.Cut(gotCode, tt.prefix); ok {
				if before, _, ok := strings.Cut(after, tt.suffix); ok {
					got = before
				} else {
					t.Fatal("suffix not found in copied file")
				}
			} else {
				t.Fatal("prefix not found in copied file")
			}
			if got != want {
				if *fixCopy {
					if err := os.WriteFile("build.go", []byte(strings.Replace(gotCode, got, want, 1)), 0o666); err != nil {
						t.Fatal(err)
					}
				} else {
					t.Error("copy is not the same as original; use '-fixcopy' to replace copied code with the code from the original file")
				}
			}
		})
	}
}
