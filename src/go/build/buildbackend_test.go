// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"reflect"
	"testing"
)

// Check that the crypto backend tag logic works and collects AllTags.
// This is based on the TestAllTags test.
func TestCryptoBackendAllTags(t *testing.T) {
	ctxt := Default
	// Remove tool tags so these tests behave the same regardless of the
	// goexperiments that happen to be set during the run.
	ctxt.ToolTags = []string{}
	ctxt.GOARCH = "amd64"
	ctxt.GOOS = "linux"
	ctxt.BuildTags = []string{"goexperiment.systemcrypto"}
	p, err := ctxt.ImportDir("testdata/backendtags_openssl", 0)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"goexperiment.opensslcrypto", "goexperiment.systemcrypto"}
	if !reflect.DeepEqual(p.AllTags, want) {
		t.Errorf("AllTags = %v, want %v", p.AllTags, want)
	}
	wantFiles := []string{"main.go", "openssl.go"}
	if !reflect.DeepEqual(p.GoFiles, wantFiles) {
		t.Errorf("GoFiles = %v, want %v", p.GoFiles, wantFiles)
	}

	ctxt.GOARCH = "amd64"
	ctxt.GOOS = "windows"
	ctxt.BuildTags = []string{"goexperiment.cngcrypto"}
	p, err = ctxt.ImportDir("testdata/backendtags_openssl", 0)
	if err != nil {
		t.Fatal(err)
	}
	// Given the current GOOS (windows), systemcrypto would not affect the
	// decision, so we don't want it to be included in AllTags.
	want = []string{"goexperiment.opensslcrypto"}
	if !reflect.DeepEqual(p.AllTags, want) {
		t.Errorf("AllTags = %v, want %v", p.AllTags, want)
	}
	wantFiles = []string{"main.go"}
	if !reflect.DeepEqual(p.GoFiles, wantFiles) {
		t.Errorf("GoFiles = %v, want %v", p.GoFiles, wantFiles)
	}

	// We want systemcrypto when cngcrypto is enabled on Windows.
	p, err = ctxt.ImportDir("testdata/backendtags_system", 0)
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"goexperiment.boringcrypto", "goexperiment.cngcrypto", "goexperiment.opensslcrypto", "goexperiment.systemcrypto"}
	if !reflect.DeepEqual(p.AllTags, want) {
		t.Errorf("AllTags = %v, want %v", p.AllTags, want)
	}
	wantFiles = []string{"main.go", "systemcrypto.go"}
	if !reflect.DeepEqual(p.GoFiles, wantFiles) {
		t.Errorf("GoFiles = %v, want %v", p.GoFiles, wantFiles)
	}
}
