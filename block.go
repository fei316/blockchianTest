package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

var genisInfo = "2019年10月6日完成，开始进行调试"

var difficulty uint64 = 23

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
	Transactions []*Transaction
}

//创建区块
func NewBlock(txs []*Transaction, prevHash []byte) *Block {

	block := Block{
		Version:    00,
		PrevHash:   prevHash,

		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: difficulty,
		Nonce:      0,
		Hash:       []byte{},
		Transactions:txs,
	}
	pow := NewProofOfWork(&block)

	hash, nonce := pow.Run()
	block.Nonce = nonce
	if !pow.IsValid() {
		log.Panic("校验工作量证明失败")
	}
	block.Hash = hash


	return &block
}

//uint转[]byte
func Unit64ToByte(num uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))
	return buf
}


func GenesisBlock(address string) *Block {
	coinBase := NewCoinbaseTx(address, genisInfo)
	block := NewBlock([]*Transaction{coinBase}, []byte{})

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
		log.Panic(err)
	}
	return block
}

func (block *Block) SetMerkelRoot() {
	var tmp [][]byte
	for _, tx := range block.Transactions {
		tmp = append(tmp, tx.TXID)
	}
	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	block.MerkelRoot = hash[:]
}