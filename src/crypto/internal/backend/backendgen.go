// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backend

// Generate files and Go code based on the set of backends:
//
// - The build constraint in nobackend.go.
// - Go files in the runtime package that detect issues with backend selection
//   and report an error at compile time.
//
// Runs in -mod=readonly mode so that it is able to run during each crypto
// backend patch. This is before the final vendoring refresh patch, so it would
// normally fail to build due to inconsistent vendoring.

// Use "go generate -run TestGenerated crypto/internal/backend"
// to run only this generator.

//go:generate go test -run TestGenerated -fix
