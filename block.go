package main

import "crypto/sha256"

//区块结构体
type Block struct {
	PrevHash []byte
	Hash     []byte
	Data     []byte
}

//创建区块
func NewBloack(data string, prevHash []byte) *Block {
	block := Block{
		PrevHash: prevHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block
}

//给区块生成hash
func (block *Block) SetHash() {
	blockInfo := append(block.PrevHash, block.Data...)
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

func GenesisBlock() *Block {
	block := NewBloack("我的创世区块2019年5月26日", []byte{})
	return block
}
