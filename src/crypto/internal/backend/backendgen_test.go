// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package backend

import (
	"bytes"
	"flag"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var fix = flag.Bool("fix", false, "if true, update the generated files to the wanted value")

const runtimePackageDir = "../../../runtime"

// backendErrPrefix is the prefix of the generated backend error files. Any file
// in the runtime package with this prefix will be considered a backend error
// file, so it's important that this prefix is unique or this generator may
// delete unexpected files.
const backendErrPrefix = "backenderr_gen_"

const generateInstruction = "run 'go generate crypto/internal/backend' to fix"

// TestGeneratedBackendErrorFiles tests that the current nobackend constraint
// is correct.
//
// Generate the build constraint in nobackend.go. This build constraint enables
// nobackend when all of the backends are not enabled. This constraint is fairly
// long and would not be trivial to maintain manually.
func TestGeneratedNobackendConstraint(t *testing.T) {
	backends := parseBackends(t)
	// none is a constraint that is met when all crypto backend constraints are
	// unmet. (That is: no backend constraint is met.)
	var none constraint.Expr
	for _, b := range backends {
		notB := &constraint.NotExpr{X: b.constraint}
		if none == nil {
			none = notB
		} else {
			none = &constraint.AndExpr{
				X: none,
				Y: notB,
			}
		}
	}
	bytes, err := os.ReadFile("nobackend.go")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(bytes), "\n")

	var gotIndex int
	var gotLine string
	for i, line := range lines {
		if strings.HasPrefix(line, "//go:build ") {
			gotIndex, gotLine = i, line
			break
		}
	}
	_ = gotIndex

	var wantLine string
	if none == nil {
		// If there are no backends yet, use a trivially true constraint.
		// We could remove the constraint line, but this would make generation
		// more complicated.
		wantLine = "//go:build go1.1"
	} else {
		wantLine = "//go:build " + none.String()
	}
	if wantLine != gotLine {
		if *fix {
			lines[gotIndex] = wantLine
			want := strings.Join(lines, "\n")
			if err := os.WriteFile("nobackend.go", []byte(want), 0o666); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Errorf("nobackend.go build constraint:\ngot %q\nwant %q\n%v", gotLine, wantLine, generateInstruction)
		}
	}
}

// TestGeneratedBackendErrorFiles tests that the current backend error files are
// the same as what would generated under the current conditions.
//
// The error files are Go files that detect issues with the backend selection
// and report an error at compile time.
//
// The issue detection files are placed in the runtime package rather than the
// crypto/internal/backend package to make sure these helpful errors will show
// up. If the files were in the backend package, DuplicateDecl and other errors
// would show up first, causing these informative errors to be skipped because
// there are too many total errors already reported. The errors would also show
// up if we put the files in the crypto package rather than the runtime package.
// (Crypto is imported before the backend backage, so the errors would show up.)
// However, then these errors would show up only if the Go program is using
// crypto. This could cause a confusing situation: if the user has a
// misconfigured backend and doesn't use crypto in their Go app, they will not
// get any errors. If they start using crypto later, they would only then get an
// error, but the cause would be much less apparent.
func TestGeneratedBackendErrorFiles(t *testing.T) {
	// Chip away at a list of files that should come from this generator.
	// Any remaining are unexpected.
	existingFiles := make(map[string]struct{})
	entries, err := os.ReadDir(runtimePackageDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), backendErrPrefix) && strings.HasSuffix(e.Name(), ".go") {
			existingFiles[filepath.Join(runtimePackageDir, e.Name())] = struct{}{}
		}
	}

	backends := parseBackends(t)
	for i := 0; i < len(backends); i++ {
		for j := i + 1; j < len(backends); j++ {
			f := testConflict(t, backends[i].name, backends[j].name)
			delete(existingFiles, f)
		}
		f := testPreventUnintendedFallback(t, backends[i])
		delete(existingFiles, f)
	}
	f := testUnsatisfied(t, backends)
	delete(existingFiles, f)
	f = testRequireFIPSWithoutBackend(t)
	delete(existingFiles, f)

	for f := range existingFiles {
		if *fix {
			if err := os.Remove(f); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Errorf("unexpected file: %q", f)
		}
	}
	if !*fix && len(existingFiles) > 0 {
		t.Log(generateInstruction)
	}
}

// testConflict checks/generates a file that fails if two backends are enabled
// at the same time.
func testConflict(t *testing.T, a, b string) string {
	return testErrorFile(
		t,
		filepath.Join(runtimePackageDir, backendErrPrefix+"conflict_"+a+"_"+b+".go"),
		"//go:build goexperiment."+a+"crypto && goexperiment."+b+"crypto",
		"The "+a+" and "+b+" backends are both enabled, but they are mutually exclusive.",
		"Please make sure only one crypto backend experiment is enabled by GOEXPERIMENT or '-tags'.")
}

func testPreventUnintendedFallback(t *testing.T, backend *backend) string {
	expTag := &constraint.TagExpr{Tag: "goexperiment." + backend.name + "crypto"}
	optOutTag := &constraint.TagExpr{Tag: "goexperiment.allowcryptofallback"}
	c := constraint.AndExpr{
		X: &constraint.AndExpr{
			X: expTag,
			Y: &constraint.NotExpr{X: backend.constraint},
		},
		Y: &constraint.NotExpr{X: optOutTag},
	}
	return testErrorFile(
		t,
		filepath.Join(runtimePackageDir, backendErrPrefix+"nofallback_"+backend.name+".go"),
		"//go:build "+c.String(),
		"The "+expTag.String()+" tag is specified, but other tags required to enable that backend were not met.",
		"Required build tags:",
		"  "+backend.constraint.String(),
		"Please check your build environment and build command for a reason one or more of these tags weren't specified.",
		"",
		"If you only performed a Go toolset upgrade and didn't expect this error, your code was likely depending on fallback to Go standard library crypto.",
		"As of Go 1.21, Go crypto fallback is a build error. This helps prevent accidental fallback.",
		"Removing "+backend.name+"crypto will restore pre-1.21 behavior by intentionally using Go standard library crypto.",
		"")
}

// testUnsatisfied checks/generates a file that fails if systemcrypto is enabled
// on an OS with no suitable backend.
func testUnsatisfied(t *testing.T, backends []*backend) string {
	constraint := "//go:build goexperiment.systemcrypto"
	for _, b := range backends {
		constraint += ` && !goexperiment.` + b.name + "crypto"
	}
	return testErrorFile(
		t,
		filepath.Join(runtimePackageDir, backendErrPrefix+"systemcrypto_nobackend.go"),
		constraint,
		"The systemcrypto feature is enabled, but it was unable to enable an appropriate crypto backend for the target GOOS.")
}

func testRequireFIPSWithoutBackend(t *testing.T) string {
	return testErrorFile(
		t,
		filepath.Join(runtimePackageDir, backendErrPrefix+"requirefips_nosystemcrypto.go"),
		"//go:build requirefips && !goexperiment.systemcrypto",
		"The requirefips tag is enabled, but no crypto backend is enabled.",
		"A crypto backend is required to enable FIPS mode.")
}

// testErrorFile checks/generates a Go file with a given build constraint that
// fails to compile. The file uses an unused string to convey an error message
// to the dev on the "go build" command line.
func testErrorFile(t *testing.T, file, constraint string, message ...string) string {
	const header = `// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is generated by crypto/internal/backend. DO NOT EDIT. DO NOT manually create files with the prefix "` + backendErrPrefix + `".`
	c := header + "\n\n" + constraint + "\n\npackage runtime\n\nfunc init() {\n\t`\n"
	for _, m := range message {
		c += "\t" + m + "\n"
	}
	c += "\tFor more information, visit https://github.com/microsoft/go/tree/microsoft/main/eng/doc/fips\n"
	c += "\t`" + "\n}\n"
	if *fix {
		if err := os.WriteFile(file, []byte(c), 0o666); err != nil {
			t.Fatal(err)
		}
	} else {
		existing, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(existing, []byte(c)) {
			t.Errorf("file %v doesn't match expected value; %v", file, generateInstruction)
			t.Log("found:", string(existing))
			t.Log("would generate:", c)
		}
	}
	return file
}

type backend struct {
	filename   string
	name       string
	constraint constraint.Expr
}

func parseBackends(t *testing.T) []*backend {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, ".", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	var bs []*backend
	// Any file in this dir that defines "Enabled" is a backend.
	for k, v := range pkgs["backend"].Files {
		if en := v.Scope.Lookup("Enabled"); en != nil {
			// nobackend defines Enabled, but it is specifically not a backend.
			if k == "nobackend.go" {
				continue
			}
			b := backend{filename: k}
			b.name, _, _ = strings.Cut(strings.TrimSuffix(k, ".go"), "_")
			for _, comment := range v.Comments {
				for _, c := range comment.List {
					if strings.HasPrefix(c.Text, "//go:build ") {
						if c, err := constraint.Parse(c.Text); err == nil {
							b.constraint = c
						} else {
							t.Fatal(err)
						}
					}
				}
			}
			bs = append(bs, &b)
		}
	}
	sort.Slice(bs, func(i, j int) bool {
		return bs[i].name < bs[j].name
	})
	return bs
}
