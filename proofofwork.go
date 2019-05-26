package main

import "math/big"

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pf := ProofOfWork{
		block:block,
	}
	targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	tmpbig := big.Int{}
	tmpbig.SetString(targetStr, 16)
	pf.target = &tmpbig
	return &pf
}

func (pow *ProofOfWork) Run() (hash []byte, nonce uint64) {
	//TODO
	return []byte("hello"), 100
}