// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build goexperiment.cngcrypto

package bbig

import "github.com/microsoft/go-crypto-winnative/cng/bbig"

var Enc = bbig.Enc
var Dec = bbig.Dec
