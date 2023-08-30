// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Do not edit the build constraint by hand. It is generated by "backendgen.go".

//go:build !(goexperiment.boringcrypto && linux && cgo && (amd64 || arm64) && !android && !msan) && !(goexperiment.cngcrypto && windows) && !(goexperiment.opensslcrypto && linux && cgo)

package backend

import (
	"crypto"
	"crypto/cipher"
	"hash"
)

const Enabled = false

type BigInt = []uint

type randReader int

func (randReader) Read(b []byte) (int, error) { panic("cryptobackend: not available") }

const RandReader = randReader(0)

func NewSHA1() hash.Hash   { panic("cryptobackend: not available") }
func NewSHA224() hash.Hash { panic("cryptobackend: not available") }
func NewSHA256() hash.Hash { panic("cryptobackend: not available") }
func NewSHA384() hash.Hash { panic("cryptobackend: not available") }
func NewSHA512() hash.Hash { panic("cryptobackend: not available") }

func NewSHA3_256() hash.Hash { panic("cryptobackend: not available") }

func SHA1(p []byte) (sum [20]byte)   { panic("cryptobackend: not available") }
func SHA224(p []byte) (sum [28]byte) { panic("cryptobackend: not available") }
func SHA256(p []byte) (sum [32]byte) { panic("cryptobackend: not available") }
func SHA384(p []byte) (sum [48]byte) { panic("cryptobackend: not available") }
func SHA512(p []byte) (sum [64]byte) { panic("cryptobackend: not available") }

func SHA3_256(p []byte) (sum [32]byte) { panic("cryptobackend: not available") }

func NewHMAC(h func() hash.Hash, key []byte) hash.Hash { panic("cryptobackend: not available") }

func NewAESCipher(key []byte) (cipher.Block, error) { panic("cryptobackend: not available") }
func NewGCMTLS(c cipher.Block) (cipher.AEAD, error) { panic("cryptobackend: not available") }

type PublicKeyECDSA struct{ _ int }
type PrivateKeyECDSA struct{ _ int }

func GenerateKeyECDSA(curve string) (X, Y, D BigInt, err error) {
	panic("cryptobackend: not available")
}
func NewPrivateKeyECDSA(curve string, X, Y, D BigInt) (*PrivateKeyECDSA, error) {
	panic("cryptobackend: not available")
}
func NewPublicKeyECDSA(curve string, X, Y BigInt) (*PublicKeyECDSA, error) {
	panic("cryptobackend: not available")
}
func SignMarshalECDSA(priv *PrivateKeyECDSA, hash []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func VerifyECDSA(pub *PublicKeyECDSA, hash []byte, sig []byte) bool {
	panic("cryptobackend: not available")
}

type PublicKeyRSA struct{ _ int }
type PrivateKeyRSA struct{ _ int }

func DecryptRSAOAEP(h, mgfHash hash.Hash, priv *PrivateKeyRSA, ciphertext, label []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func DecryptRSAPKCS1(priv *PrivateKeyRSA, ciphertext []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func DecryptRSANoPadding(priv *PrivateKeyRSA, ciphertext []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func EncryptRSAOAEP(h, mgfHash hash.Hash, pub *PublicKeyRSA, msg, label []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func EncryptRSAPKCS1(pub *PublicKeyRSA, msg []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func EncryptRSANoPadding(pub *PublicKeyRSA, msg []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func GenerateKeyRSA(bits int) (N, E, D, P, Q, Dp, Dq, Qinv BigInt, err error) {
	panic("cryptobackend: not available")
}
func NewPrivateKeyRSA(N, E, D, P, Q, Dp, Dq, Qinv BigInt) (*PrivateKeyRSA, error) {
	panic("cryptobackend: not available")
}
func NewPublicKeyRSA(N, E BigInt) (*PublicKeyRSA, error) {
	panic("cryptobackend: not available")
}
func SignRSAPKCS1v15(priv *PrivateKeyRSA, h crypto.Hash, hashed []byte) ([]byte, error) {
	panic("cryptobackend: not available")
}
func SignRSAPSS(priv *PrivateKeyRSA, h crypto.Hash, hashed []byte, saltLen int) ([]byte, error) {
	panic("cryptobackend: not available")
}
func VerifyRSAPKCS1v15(pub *PublicKeyRSA, h crypto.Hash, hashed, sig []byte) error {
	panic("cryptobackend: not available")
}
func VerifyRSAPSS(pub *PublicKeyRSA, h crypto.Hash, hashed, sig []byte, saltLen int) error {
	panic("cryptobackend: not available")
}

type PublicKeyECDH struct{}
type PrivateKeyECDH struct{}

func ECDH(*PrivateKeyECDH, *PublicKeyECDH) ([]byte, error)    { panic("cryptobackend: not available") }
func GenerateKeyECDH(string) (*PrivateKeyECDH, []byte, error) { panic("cryptobackend: not available") }
func NewPrivateKeyECDH(string, []byte) (*PrivateKeyECDH, error) {
	panic("cryptobackend: not available")
}
func NewPublicKeyECDH(string, []byte) (*PublicKeyECDH, error) { panic("cryptobackend: not available") }
func (*PublicKeyECDH) Bytes() []byte                          { panic("cryptobackend: not available") }
func (*PrivateKeyECDH) PublicKey() (*PublicKeyECDH, error)    { panic("cryptobackend: not available") }
