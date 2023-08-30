// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build goexperiment.opensslcrypto && linux && cgo

// Package openssl provides access to OpenSSLCrypto implementation functions.
// Check the variable Enabled to find out whether OpenSSLCrypto is available.
// If OpenSSLCrypto is not available, the functions in this package all panic.
package backend

import (
	"hash"

	_ "crypto/boring"
	_ "unsafe"
)

// Enabled controls whether FIPS crypto is enabled.
const Enabled = true

//go:linkname RandReader crypto/internal/backend.RandReader
var RandReader int

// func NewSHA1() hash.Hash
// func NewSHA224() hash.Hash
// func NewSHA256() hash.Hash
// func NewSHA384() hash.Hash
// func NewSHA512() hash.Hash

//go:linkname NewSHA3_256 crypto/internal/backend.NewSHA3_256
func NewSHA3_256() hash.Hash

// func SHA1(p []byte) (sum [20]byte)   { return openssl.SHA1(p) }
// func SHA224(p []byte) (sum [28]byte) { return openssl.SHA224(p) }
// func SHA256(p []byte) (sum [32]byte) { return openssl.SHA256(p) }
// func SHA384(p []byte) (sum [48]byte) { return openssl.SHA384(p) }
// func SHA512(p []byte) (sum [64]byte) { return openssl.SHA512(p) }

//go:linkname SHA3_256 crypto/internal/backend.SHA3_256
func SHA3_256(p []byte) (sum [32]byte)

// func NewHMAC(h func() hash.Hash, key []byte) hash.Hash { return openssl.NewHMAC(h, key) }

// func NewAESCipher(key []byte) (cipher.Block, error) { return openssl.NewAESCipher(key) }
// func NewGCMTLS(c cipher.Block) (cipher.AEAD, error) { return openssl.NewGCMTLS(c) }

// type PublicKeyECDSA = openssl.PublicKeyECDSA
// type PrivateKeyECDSA = openssl.PrivateKeyECDSA

// func GenerateKeyECDSA(curve string) (X, Y, D openssl.BigInt, err error) {
// 	return openssl.GenerateKeyECDSA(curve)
// }

// func NewPrivateKeyECDSA(curve string, X, Y, D openssl.BigInt) (*openssl.PrivateKeyECDSA, error) {
// 	return openssl.NewPrivateKeyECDSA(curve, X, Y, D)
// }

// func NewPublicKeyECDSA(curve string, X, Y openssl.BigInt) (*openssl.PublicKeyECDSA, error) {
// 	return openssl.NewPublicKeyECDSA(curve, X, Y)
// }

// func SignMarshalECDSA(priv *openssl.PrivateKeyECDSA, hash []byte) ([]byte, error) {
// 	return openssl.SignMarshalECDSA(priv, hash)
// }

// func VerifyECDSA(pub *openssl.PublicKeyECDSA, hash []byte, sig []byte) bool {
// 	return openssl.VerifyECDSA(pub, hash, sig)
// }

// type PublicKeyRSA = openssl.PublicKeyRSA
// type PrivateKeyRSA = openssl.PrivateKeyRSA

// func DecryptRSAOAEP(h, mgfHash hash.Hash, priv *openssl.PrivateKeyRSA, ciphertext, label []byte) ([]byte, error) {
// 	return openssl.DecryptRSAOAEP(h, mgfHash, priv, ciphertext, label)
// }

// func DecryptRSAPKCS1(priv *openssl.PrivateKeyRSA, ciphertext []byte) ([]byte, error) {
// 	return openssl.DecryptRSAPKCS1(priv, ciphertext)
// }

// func DecryptRSANoPadding(priv *openssl.PrivateKeyRSA, ciphertext []byte) ([]byte, error) {
// 	return openssl.DecryptRSANoPadding(priv, ciphertext)
// }

// func EncryptRSAOAEP(h, mgfHash hash.Hash, pub *openssl.PublicKeyRSA, msg, label []byte) ([]byte, error) {
// 	return openssl.EncryptRSAOAEP(h, mgfHash, pub, msg, label)
// }

// func EncryptRSAPKCS1(pub *openssl.PublicKeyRSA, msg []byte) ([]byte, error) {
// 	return openssl.EncryptRSAPKCS1(pub, msg)
// }

// func EncryptRSANoPadding(pub *openssl.PublicKeyRSA, msg []byte) ([]byte, error) {
// 	return openssl.EncryptRSANoPadding(pub, msg)
// }

// func GenerateKeyRSA(bits int) (N, E, D, P, Q, Dp, Dq, Qinv openssl.BigInt, err error) {
// 	return openssl.GenerateKeyRSA(bits)
// }

// func NewPrivateKeyRSA(N, E, D, P, Q, Dp, Dq, Qinv openssl.BigInt) (*openssl.PrivateKeyRSA, error) {
// 	return openssl.NewPrivateKeyRSA(N, E, D, P, Q, Dp, Dq, Qinv)
// }

// func NewPublicKeyRSA(N, E openssl.BigInt) (*openssl.PublicKeyRSA, error) {
// 	return openssl.NewPublicKeyRSA(N, E)
// }

// func SignRSAPKCS1v15(priv *openssl.PrivateKeyRSA, h crypto.Hash, hashed []byte) ([]byte, error) {
// 	return openssl.SignRSAPKCS1v15(priv, h, hashed)
// }

// func SignRSAPSS(priv *openssl.PrivateKeyRSA, h crypto.Hash, hashed []byte, saltLen int) ([]byte, error) {
// 	return openssl.SignRSAPSS(priv, h, hashed, saltLen)
// }

// func VerifyRSAPKCS1v15(pub *openssl.PublicKeyRSA, h crypto.Hash, hashed, sig []byte) error {
// 	return openssl.VerifyRSAPKCS1v15(pub, h, hashed, sig)
// }

// func VerifyRSAPSS(pub *openssl.PublicKeyRSA, h crypto.Hash, hashed, sig []byte, saltLen int) error {
// 	return openssl.VerifyRSAPSS(pub, h, hashed, sig, saltLen)
// }

// type PublicKeyECDH = openssl.PublicKeyECDH
// type PrivateKeyECDH = openssl.PrivateKeyECDH

// func ECDH(priv *openssl.PrivateKeyECDH, pub *openssl.PublicKeyECDH) ([]byte, error) {
// 	return openssl.ECDH(priv, pub)
// }

// func GenerateKeyECDH(curve string) (*openssl.PrivateKeyECDH, []byte, error) {
// 	return openssl.GenerateKeyECDH(curve)
// }

// func NewPrivateKeyECDH(curve string, bytes []byte) (*openssl.PrivateKeyECDH, error) {
// 	return openssl.NewPrivateKeyECDH(curve, bytes)
// }

// func NewPublicKeyECDH(curve string, bytes []byte) (*openssl.PublicKeyECDH, error) {
// 	return openssl.NewPublicKeyECDH(curve, bytes)
// }
