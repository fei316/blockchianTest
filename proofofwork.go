package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)



type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pf := ProofOfWork{
		block: block,
	}
	targetLocal := big.NewInt(1)
	targetLocal.Lsh(targetLocal, 256 - uint(block.Difficulty))
	pf.target = targetLocal
	return &pf
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	log.Println("开始挖矿...")
	var nonce uint64
	var hash [32]byte
	for {

		hash = sha256.Sum256(pow.PrepareData(nonce))

		tempbig := new(big.Int)
		tempbig.SetBytes(hash[:])

		if tempbig.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功\nHash:%x\nNonce:%d\n", hash, nonce)

			break
		} else {
			nonce++
		}

	}
	log.Println("挖矿结束...")
	return hash[:], nonce
}

func (pow *ProofOfWork) PrepareData(nonce uint64) []byte {

	block := pow.block
	block.SetMerkelRoot()
	tmp := [][]byte{
		Unit64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Unit64ToByte(block.TimeStamp),
		Unit64ToByte(block.Difficulty),
		Unit64ToByte(nonce),

	}

	blockInfo := bytes.Join(tmp, []byte{})

	return blockInfo
}

//校验工作量证明是否为true
func (pow *ProofOfWork) IsValid() bool{
	hash := sha256.Sum256(pow.PrepareData(pow.block.Nonce))
	tempbig := new(big.Int)
	tempbig.SetBytes(hash[:])

	if tempbig.Cmp(pow.target) == -1 {
		return true
	}
	return false
}
