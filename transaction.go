package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"
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
	input := TXInput{nil, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetID()
	return &tx
}

//创建交易
func NewTransaction(from, to string, amount float64, bc *BlockChian) *Transaction {
	inputs, total := bc.FindSuitableUTXOs(from, amount)
	if total < amount {
		log.Printf("您的余额为%f，请先挣点钱再来")
		os.Exit(1)
	}
	var tran = Transaction{
		TXID:     nil,
		TXInputs: inputs,
	}

	var outputs []TXOutput
	output := TXOutput{
		value:amount,
		PubkeyHash:to,
	}
	outputs = append(outputs, output)
	if total > amount {
		zhaoling := TXOutput{
			value:total - amount,
			PubkeyHash:from,
		}
		outputs = append(outputs, zhaoling)
	}
	tran.TXOutputs = outputs

	return &tran
}



func (bc *BlockChian) FindUTXOTransactions() []Transaction {
	txo := make(map[string][]int64)
	var transcations []Transaction
	//循环bc
	bcInterator := bc.NewBlockchainInterator()
	block := bcInterator.Next()
	for {
		//循环交易
		trans := block.Transactions
		TRANS:
		for i:=0; i< len(trans); i++ {
			tran := trans[i]
			//循环output
			outputs := tran.TXOutputs
			OUTPUTS:
			for outindex, _ := range outputs {
				if txo[string(tran.TXID)] != nil {
					indexs := txo[string(tran.TXID)]
					for _, index := range indexs {
						if index == int64(outindex) {
							continue OUTPUTS
						}
					}
				}

				//output没有被消耗，判断是否属于这个地址的
				//if output.OutputCanBeUnlocked(addr) {
					transcations = append(transcations, *tran)
					continue TRANS
				//}
			}



			if !tran.IsCoinbaseTran() {
				//循环input
				inputs := tran.TXInputs
				for _, input := range inputs {
					//if input.InputCanUnlock(addr) {
						txo[string(input.TXID)] = append(txo[string(input.TXID)], input.Index)
					//}
				}
			}

		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return transcations

}

func (bc *BlockChian) getUTXOs(address string) []TXOutput {
	var outs []TXOutput
	txs := bc.FindUTXOTransactions()
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if output.OutputCanBeUnlocked(address) {
				outs = append(outs, output)
			}
		}
	}
	return outs
}

//判断是否为coinbase交易
func (tran *Transaction) IsCoinbaseTran() bool {
	if len(tran.TXInputs) == 1 {
		if tran.TXInputs[0].TXID == nil && tran.TXInputs[0].Index == -1 {
			return true
		}
	}
	return false
}

func (output *TXOutput) OutputCanBeUnlocked(address string) bool {
	return output.PubkeyHash == address
}

func (input *TXInput) InputCanUnlock(address string) bool{
	return input.Sig == address
}