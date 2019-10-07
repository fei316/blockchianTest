package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	_ "github.com/btcsuite/btcutil/base58"
	"log"
	"math/big"
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

	Signature []byte
	PublicKey []byte
}

type TXOutput struct {
	Value float64

	PubkeyHash []byte
}

//创建output
func NewTXOutput(value float64, address string) *TXOutput {
	var output TXOutput
	output.Value = value
	output.lock(address)

	return &output
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
	input := TXInput{nil, -1, []byte(data), nil}

	output := NewTXOutput(reward, address)
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	tx.SetID()
	return &tx
}

//创建交易
func NewTransaction(from, to string, amount float64, bc *BlockChian) *Transaction {

	inputs, total, privateKey := bc.FindSuitableUTXOs(from, amount)

	if total < amount {
		log.Printf("您的余额为%f，请先挣点钱再来", total)
		os.Exit(1)
	}
	var tran = Transaction{
		TXID:     []byte{},
		TXInputs: inputs,
	}

	var outputs []TXOutput
	output := TXOutput{
		Value:amount,
	}
	output.lock(to)
	outputs = append(outputs, output)
	if total > amount {
		zhaoling := TXOutput{
			Value:total - amount,
		}
		zhaoling.lock(from)
		outputs = append(outputs, zhaoling)
	}
	tran.TXOutputs = outputs

	tran.SetID()

	var prevTrans = make(map[string]Transaction)

	for _, input := range inputs {

		tempTran, err := bc.getTransactionByID(input.TXID)

		if err != nil {
			log.Panic(err)
		}

		prevTrans[string(input.TXID)] = *tempTran

	}

	tran.sign(&privateKey, prevTrans)

	return &tran
}



func (blockchain *BlockChian) FindUTXOTransactions() []Transaction {
	txo := make(map[string][]int64)
	var transcations []Transaction
	//循环bc
	bcInterator := blockchain.NewBlockchainInterator()
	for {

		block := bcInterator.Next()

		//循环交易
		trans := block.Transactions
		TRANS:
		for i:=0; i< len(trans); i++ {
			tran := trans[i]
			//循环output
			outputs := tran.TXOutputs


			for outindex, _ := range outputs {

				if txo[string(tran.TXID)] != nil {
					indexs := txo[string(tran.TXID)]

					for _, index := range indexs {


						if index != int64(outindex) {

							transcations = append(transcations, *tran)
							if !tran.IsCoinbaseTran() {

								//循环input
								inputs := tran.TXInputs
								for _, input := range inputs {

									txo[string(input.TXID)] = append(txo[string(input.TXID)], input.Index)
								}
							}
							continue TRANS
						}
					}
				} else {

					transcations = append(transcations, *tran)
					if !tran.IsCoinbaseTran() {

						//循环input
						inputs := tran.TXInputs
						for _, input := range inputs {

							txo[string(input.TXID)] = append(txo[string(input.TXID)], input.Index)
						}
					}
					continue TRANS
				}
			}


		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return transcations

}

func (blockchain *BlockChian) getUTXOs(pubHash []byte) []TXOutput {
	var outs []TXOutput
	txs := blockchain.FindUTXOTransactions()

	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if output.OutputCanBeUnlocked(pubHash) {
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

func (output *TXOutput) OutputCanBeUnlocked(pubHash []byte) bool {
	return bytes.Equal(output.PubkeyHash, pubHash)
}

func (input *TXInput) InputCanUnlock(pubHash []byte) bool{
	return bytes.Equal(input.Signature, pubHash)
}

//锁定
func (output *TXOutput)lock(address string)  {

	output.PubkeyHash = GetPubHashByAddress(address)
}

//签名
func (tx *Transaction)sign(privateKey *ecdsa.PrivateKey, prevTrans map[string]Transaction)  {
	//1，辅助一份交易对象
	txCopy := tx.copyTran()
	//2，给对象里的input的pubkey复制output里的pubkeyhash
	for i, input := range txCopy.TXInputs {
		txCopy.TXInputs[i].PublicKey = prevTrans[string(input.TXID)].TXOutputs[input.Index].PubkeyHash
		//3，把交易sethash
		txCopy.SetID()
		txCopy.TXInputs[i].PublicKey = nil
		//4，然后对hash数据进行签名
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, txCopy.TXID)
		if err != nil {
			log.Panic(err)
		}
		tx.TXInputs[i].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

//校验
func (tx *Transaction)verify(prevTrans map[string]Transaction) bool {
	txCopy := tx.copyTran()
	for i, input := range txCopy.TXInputs {
		pubkeyHash := prevTrans[string(input.TXID)].TXOutputs[input.Index].PubkeyHash
		txCopy.TXInputs[i].PublicKey = prevTrans[string(input.TXID)].TXOutputs[input.Index].PubkeyHash
		txCopy.TXInputs[i].Signature = nil
		//3，把交易sethash
		txCopy.SetID()
		txCopy.TXInputs[i].PublicKey = nil

		//获得公钥
		curve := elliptic.P256()
		pubkey := tx.TXInputs[i].PublicKey
		pubkeyHashTemp := PubKeyToPubHash(pubkey)
		if !bytes.Equal(pubkeyHash, pubkeyHashTemp) {
			log.Panic("input里存储的pubkey和所引用的pubkeyhash不一致")
		}
		pubkeyLen := len(pubkey)
		var x = big.Int{}
		var y = big.Int{}
		x.SetBytes(pubkey[:pubkeyLen/2])
		y.SetBytes(pubkey[pubkeyLen/2:])
		rawPubkey := ecdsa.PublicKey{curve, &x, &y}

		//获得签名
		r := big.Int{}
		s := big.Int{}
		sig := tx.TXInputs[i].Signature
		sigLen := len(sig)
		r.SetBytes(sig[:sigLen/2])
		s.SetBytes(sig[sigLen/2:])
		if !ecdsa.Verify(&rawPubkey, txCopy.TXID, &r, &s) {
			return false
		}

	}
	return true
}

//复制交易对象
func (tx *Transaction)copyTran() Transaction  {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.TXID, input.Index, nil, nil})
	}
	for _, output := range outputs {
		outputs = append(outputs, output)
	}
	return Transaction{tx.TXID, inputs, outputs}
}