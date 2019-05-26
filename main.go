package main

import (
	"crypto/sha256"
	"fmt"
)

//区块链结构体
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

func (block *Block) SetHash() {
	blockInfo := append(block.PrevHash, block.Data...)
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

func main() {
	block := NewBloack("小明向小红转了10个比特币", []byte{})
	fmt.Printf("前Hash：%x\n", block.PrevHash)
	fmt.Printf("当前Hash：%x\n", block.Hash)
	fmt.Printf("数据：%s\n", block.Data)
}
