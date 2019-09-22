package main

import (

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
	_ "github.com/btcsuite/btcutil/base58"


)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func NewWallet() *Wallet  {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKeyOri := privateKey.PublicKey
	publickKey := append(pubKeyOri.X.Bytes(), pubKeyOri.Y.Bytes()...)
	wallet := Wallet{
		PrivateKey:*privateKey,
		PublicKey:publickKey,
	}
	return &wallet
}

//获取地址
func (wallet *Wallet) getAddress() []byte {
	publicKey := wallet.PublicKey
	pubhash := PubKeyToPubHash(publicKey)
	payload := append([]byte{00}, pubhash...)
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	payload = append(payload, hash2[:4]...)
	address := base58.Encode(payload)
	return []byte(address)
}

func PubKeyToPubHash(publicKey []byte) []byte {

	pubSha256 := sha256.Sum256(publicKey)
	hasher160 := ripemd160.New()
	hasher160.Write(pubSha256[:])
	pubhash := hasher160.Sum(nil)
	return pubhash
}