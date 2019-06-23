package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

var reward = 12.5

type Transaction struct {
	TXID []byte
	TXInputs []TXInput
	TXOutputs []TXOutput

}

type TXInput struct {
	TXID []byte
	Index int64
	Sig string
}

type TXOutput struct {
	value float64
	PubkeyHash string
}

//设置交易id方法
func (tx *Transaction) SetID()  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic("编码交易失败")
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]

}

//创建挖矿交易
func NewCoinbaseTx(address string, data string) *Transaction {
	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetID()
	return &tx
}

func (bc *BlockChian) getUTXOs(address string) []TXOutput {
	//TODO
	return []TXOutput{}
}