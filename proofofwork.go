package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pf := ProofOfWork{
		block:block,
	}
	targetStr := "0000000100000000000000000000000000000000000000000000000000000000"
	tmpbig := big.Int{}
	tmpbig.SetString(targetStr, 16)
	pf.target = &tmpbig
	return &pf
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	block := pow.block
	var nonce uint64
	var hash [32]byte
	for {
		tmp := [][]byte{
			Unit64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Unit64ToByte(block.TimeStamp),
			Unit64ToByte(block.Difficulty),
			Unit64ToByte(nonce),
			block.Data,
		}

		blockInfo := bytes.Join(tmp, []byte{})

		hash = sha256.Sum256(blockInfo)

		tempbig := big.Int{}
		tempbig.SetBytes(hash[:])

		if tempbig.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功\nHash:%x\nNonce:%d\n", hash, nonce)
			break
		} else {
			nonce ++
		}

	}


	return hash[:], nonce
}