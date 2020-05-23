package bsign

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"math/big"
)

// Sign is signing ballot provided by client.
//
// Function takes key and ballot. It returns ballot^d mod N.
func Sign(key *rsa.PrivateKey, in []byte) *big.Int {
	// Calculate m^d to generate sign for client.
	m := new(big.Int).SetBytes(in)
	return m.Exp(m, key.D, key.PublicKey.N)
}

// Verify check if the sign is Valid.
//
// Input is a public key corresponding to key used for signing and sign, which is a pair
// of numbers mod N. Sign is valid if hash(m) = (md)^e mod N.
func Verify(key *rsa.PublicKey, m []byte, md []byte) bool {
	// We check if hash(m) = (md)^e mod N.
	// If sign is correct, md = hash(m)^d and equality is satisfied.
	hash := sha256.Sum256(m)
	// Calculate md^e mod N
	mdi := new(big.Int).SetBytes(md)
	bhi := new(big.Int).Exp(mdi, big.NewInt(int64(key.E)), key.N)

	if subtle.ConstantTimeCompare(bhi.Bytes(), hash[:]) >= 1 {
		return true
	}
	return false
}
