// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !boringcrypto && !goexperiment.opensslcrypto && !goexperiment.cngcrypto

package ecdsa

import boring "crypto/internal/backend"

func boringPublicKey(*PublicKey) (*boring.PublicKeyECDSA, error) {
	panic("boringcrypto: not available")
}
func boringPrivateKey(*PrivateKey) (*boring.PrivateKeyECDSA, error) {
	panic("boringcrypto: not available")
}
