package main

import (
	"crypto/sha256"
	"fmt"
)

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

//区块链结构体
type BlockChian struct {
	blocks []Block
}

func GenesisBlock() *Block {
	block := NewBloack("我的创世区块2019年5月26日", []byte{})
	return block
}

//创建区块链
func NewBlockchian() *BlockChian {
	blockchian := BlockChian{
		blocks:[]Block{*GenesisBlock()},
	}
	return &blockchian
}

func main() {
	bc := NewBlockchian()
	for i, block := range bc.blocks {
		fmt.Printf("=====当前区块高度：%d =======\n", i)
		fmt.Printf("前Hash：%x\n", block.PrevHash)
		fmt.Printf("当前Hash：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", block.Data)
	}

}
