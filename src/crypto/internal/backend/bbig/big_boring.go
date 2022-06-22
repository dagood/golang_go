// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build goexperiment.boringcrypto

package bbig

import "crypto/internal/boring/bbig"

var Enc = bbig.Enc
var Dec = bbig.Dec
