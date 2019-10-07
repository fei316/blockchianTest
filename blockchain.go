package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

const blockchaindb = "blockchain.db"
const blockbucket = "blockbucket"

//区块链结构体
type BlockChian struct {
	db   *bolt.DB
	Tail []byte
}

//创建区块链
func NewBlockchian() *BlockChian {

	db, err := bolt.Open(blockchaindb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()
	if err != nil {
		log.Panic("open db err")
	}

	var tail = []byte{}

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockbucket))
			if err != nil {
				log.Panic("create bucket err")
			}
			ws := NewWallets()
			address := ws.createWallet()
			block := GenesisBlock(address)
			err := bucket.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put([]byte("lastHash"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			tail = block.Hash

		}

		tail = bucket.Get([]byte("lastHash"))

		return nil
	})
	return &BlockChian{
		db:   db,
		Tail: tail,
	}
}

//区块链添加区块
func (blockchain *BlockChian) AddBlock(txs []*Transaction) {

	for _, tx := range txs {

		if !blockchain.verifyTrans(tx) {
			log.Print("交易校验失败，终止交易")
			return
		}
	}

	db := blockchain.db
	tail := blockchain.Tail
	db, err := bolt.Open(blockchaindb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if bucket == nil {
			log.Panic("bucket为空，不应该为空")
		}
		block := NewBlock(txs, tail)
		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte("lastHash"), block.Hash)
		if err != nil {
			log.Panic(err)
		}
		blockchain.Tail = block.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}


}

//获取某地址下足够多钱多utxos
func (blockchain *BlockChian) FindSuitableUTXOs(address string, amount float64) ([]TXInput, float64, ecdsa.PrivateKey){
	ws := NewWallets()
	wallet := ws.WalletsMap[address]
	publicKey := wallet.PublicKey
	privateKey := wallet.PrivateKey
	txo := make(map[string][]int64)
	var utxos []TXInput
	var total float64 = 0
	//循环bc
	bcInterator := blockchain.NewBlockchainInterator()

	BLOCK:
	for {

		block := bcInterator.Next()
		//循环交易
		trans := block.Transactions

		for i:=0; i< len(trans); i++ {
			tran := trans[i]
			//循环output
			outputs := tran.TXOutputs
		OUTPUTS:
			for outindex, output := range outputs {
				if txo[string(tran.TXID)] != nil {
					indexs := txo[string(tran.TXID)]
					for _, index := range indexs {
						if index == int64(outindex) {
							continue OUTPUTS
						}
					}
				}

				//output没有被消耗，判断是否属于这个地址的
				if output.OutputCanBeUnlocked(PubKeyToPubHash(publicKey)) {
					tmpinput := TXInput{
						TXID:tran.TXID,
						Index:int64(outindex),
						Signature:nil,
						PublicKey:publicKey,
					}
					if total < amount{
						utxos = append(utxos, tmpinput)
						total += output.Value
					} else {
						break BLOCK
					}

				}
			}



			if !tran.IsCoinbaseTran() {
				//循环input
				inputs := tran.TXInputs
				for _, input := range inputs {
					if input.InputCanUnlock(PubKeyToPubHash(publicKey)) {
						txo[string(input.TXID)] = append(txo[string(input.TXID)], input.Index)
					}
				}
			}

		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return utxos, total, privateKey
}

//获取区块链
func GetBlockchian() *BlockChian {
	db, err := bolt.Open(blockchaindb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()
	if err != nil {
		log.Panic(err)

	}
	var tail []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if bucket == nil {
			log.Panic("create bucket err")
		}

		tail = bucket.Get([]byte("lastHash"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := BlockChian{
		db:db,
		Tail:tail,
	}

	return &bc
}

//根据ID获取交易
func (blockchain *BlockChian)getTransactionByID(id []byte) (*Transaction, error) {

	var transaction *Transaction
	var flag bool = false
	//循环bc
	interator := blockchain.NewBlockchainInterator()

	for {
		block := interator.Next()
		//循环交易
		trans := block.Transactions
		for _, tran := range trans {
			if bytes.Equal(tran.TXID, id) {
				transaction = tran
				flag = true
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	if flag {
		return transaction, nil
	} else {
		return transaction, errors.New("根据ID没有找到交易")
	}
}

func (blockchain *BlockChian)verifyTrans(tx *Transaction) bool {
	if tx.IsCoinbaseTran() {

		return true
	}
	var prevTrans = make(map[string]Transaction)

	for _, input := range tx.TXInputs {
		tempTran, err := blockchain.getTransactionByID(input.TXID)
		if err != nil {
			log.Panic(err)
		}

		prevTrans[string(input.TXID)] = *tempTran

	}
	return tx.verify(prevTrans)
}
