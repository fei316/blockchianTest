package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

//区块结构体
type Block struct {
	Version uint64

	PrevHash []byte

	MerkelRoot []byte

	TimeStamp uint64

	Difficulty uint64

	Nonce uint64

	//hash和data，正常情况布不存储在这里，存储在这里是为了实现方便
	Hash []byte
	Data []byte
}

//创建区块
func NewBloack(data string, prevHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//uint转[]byte
func Unit64ToByte(num uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))
	return buf
}

//给区块生成hash
func (block *Block) SetHash() {

	tmp := [][]byte{
		Unit64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Unit64ToByte(block.TimeStamp),
		Unit64ToByte(block.Difficulty),
		Unit64ToByte(block.Nonce),
		block.Data,
	}

	blockInfo := bytes.Join(tmp, []byte{})

	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

func GenesisBlock() *Block {
	block := NewBloack("我的创世区块2019年5月26日", []byte{})
	return block
}

//序列化区块
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("序列化区块失败")
	}
	return buffer.Bytes()
}

//序列化区块
func Serialize(block *Block) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("序列化区块失败")
	}
	return buffer.Bytes()
}

//反序列化区块
func DeSerialize(data []byte) Block {

	decoder := gob.NewDecoder(bytes.NewReader(data))
	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("反序列化区块失败")
	}
	return block
}
