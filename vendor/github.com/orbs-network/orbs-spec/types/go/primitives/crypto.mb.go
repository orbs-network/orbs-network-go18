// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package primitives

import (
	"bytes"
	"fmt"
)

type NodeAddress []byte

func (x NodeAddress) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x NodeAddress) Equal(y NodeAddress) bool {
	return bytes.Equal(x, y)
}

func (x NodeAddress) KeyForMap() string {
	return string(x)
}

type Sha256 []byte

func (x Sha256) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Sha256) Equal(y Sha256) bool {
	return bytes.Equal(x, y)
}

func (x Sha256) KeyForMap() string {
	return string(x)
}

type MerkleSha256 []byte

func (x MerkleSha256) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x MerkleSha256) Equal(y MerkleSha256) bool {
	return bytes.Equal(x, y)
}

func (x MerkleSha256) KeyForMap() string {
	return string(x)
}

type Keccak256 []byte

func (x Keccak256) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Keccak256) Equal(y Keccak256) bool {
	return bytes.Equal(x, y)
}

func (x Keccak256) KeyForMap() string {
	return string(x)
}

type Ed25519Sig []byte

func (x Ed25519Sig) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Ed25519Sig) Equal(y Ed25519Sig) bool {
	return bytes.Equal(x, y)
}

func (x Ed25519Sig) KeyForMap() string {
	return string(x)
}

type Ed25519PublicKey []byte

func (x Ed25519PublicKey) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Ed25519PublicKey) Equal(y Ed25519PublicKey) bool {
	return bytes.Equal(x, y)
}

func (x Ed25519PublicKey) KeyForMap() string {
	return string(x)
}

type Ed25519PrivateKey []byte

func (x Ed25519PrivateKey) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Ed25519PrivateKey) Equal(y Ed25519PrivateKey) bool {
	return bytes.Equal(x, y)
}

func (x Ed25519PrivateKey) KeyForMap() string {
	return string(x)
}

type Ripemd160Sha256 []byte

func (x Ripemd160Sha256) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Ripemd160Sha256) Equal(y Ripemd160Sha256) bool {
	return bytes.Equal(x, y)
}

func (x Ripemd160Sha256) KeyForMap() string {
	return string(x)
}

type Bls1Sig []byte

func (x Bls1Sig) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Bls1Sig) Equal(y Bls1Sig) bool {
	return bytes.Equal(x, y)
}

func (x Bls1Sig) KeyForMap() string {
	return string(x)
}

type Bls1PublicKey []byte

func (x Bls1PublicKey) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Bls1PublicKey) Equal(y Bls1PublicKey) bool {
	return bytes.Equal(x, y)
}

func (x Bls1PublicKey) KeyForMap() string {
	return string(x)
}

type Bls1PrivateKey []byte

func (x Bls1PrivateKey) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x Bls1PrivateKey) Equal(y Bls1PrivateKey) bool {
	return bytes.Equal(x, y)
}

func (x Bls1PrivateKey) KeyForMap() string {
	return string(x)
}

type EcdsaSecp256K1Sig []byte

func (x EcdsaSecp256K1Sig) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x EcdsaSecp256K1Sig) Equal(y EcdsaSecp256K1Sig) bool {
	return bytes.Equal(x, y)
}

func (x EcdsaSecp256K1Sig) KeyForMap() string {
	return string(x)
}

type EcdsaSecp256K1PublicKey []byte

func (x EcdsaSecp256K1PublicKey) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x EcdsaSecp256K1PublicKey) Equal(y EcdsaSecp256K1PublicKey) bool {
	return bytes.Equal(x, y)
}

func (x EcdsaSecp256K1PublicKey) KeyForMap() string {
	return string(x)
}

type EcdsaSecp256K1PrivateKey []byte

func (x EcdsaSecp256K1PrivateKey) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x EcdsaSecp256K1PrivateKey) Equal(y EcdsaSecp256K1PrivateKey) bool {
	return bytes.Equal(x, y)
}

func (x EcdsaSecp256K1PrivateKey) KeyForMap() string {
	return string(x)
}

type BloomFilter []byte

func (x BloomFilter) String() string {
	return fmt.Sprintf("%x", []byte(x))
}

func (x BloomFilter) Equal(y BloomFilter) bool {
	return bytes.Equal(x, y)
}

func (x BloomFilter) KeyForMap() string {
	return string(x)
}
