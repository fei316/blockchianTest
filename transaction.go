package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Transaction struct {
	TXID []byte
	TXInputs []TXInput
	TXOutputs []TXOutput

}

type TXInput struct {
	TXID []byte
	Index int64
	sig string
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