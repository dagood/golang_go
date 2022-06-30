// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build goexperiment.opensslcrypto

package bbig

import "github.com/golang-fips/openssl/v2/bbig"

var Enc = bbig.Enc
var Dec = bbig.Dec
