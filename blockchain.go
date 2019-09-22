package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const blockchaindb = "blockchain.db"
const blockbucket = "blockbucket"

//区块链结构体
type BlockChian struct {
	db   *bolt.DB
	tail []byte
}

//创建区块链
func NewBlockchian(address string) *BlockChian {

	db, err := bolt.Open(blockchaindb, 0600, nil)

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
			block := GenesisBlock(address)
			bucket.Put(block.Hash, block.Serialize())
			bucket.Put([]byte("lastHash"), block.Hash)
			tail = block.Hash

		}

		tail = bucket.Get([]byte("lastHash"))

		return nil
	})
	return &BlockChian{
		db:   db,
		tail: tail,
	}
}

//区块链添加区块
func (blockchain *BlockChian) AddBlock(txs []*Transaction) {
	db := blockchain.db
	tail := blockchain.tail
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if bucket == nil {
			log.Panic("bucket为空，不应该为空")
		}
		block := NewBloack(txs, tail)
		bucket.Put(block.Hash, block.Serialize())

		bucket.Put([]byte("lastHash"), block.Hash)
		blockchain.tail = block.Hash
		return nil
	})


}

//获取某地址下足够多钱多utxos
func (bc *BlockChian) FindSuitableUTXOs(address string, amount float64) ([]TXInput, float64){
	ws := NewWallets()
	wallet := ws.WalletsMap[address]
	publicKey := wallet.PublicKey
	//TODO privateKey := wallet.PrivateKey
	txo := make(map[string][]int64)
	var utxos []TXInput
	var total float64 = 0
	//循环bc
	bcInterator := bc.NewBlockchainInterator()
	block := bcInterator.Next()

	BLOCK:
	for {
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
				if output.OutputCanBeUnlocked(address) {
					tmpinput := TXInput{
						TXID:tran.TXID,
						Index:int64(outindex),
						Signature:nil,
						PublicKey:publicKey,
					}
					if total < amount{
						utxos = append(utxos, tmpinput)
						total += output.value
					} else {
						break BLOCK
					}

				}
			}



			if !tran.IsCoinbaseTran() {
				//循环input
				inputs := tran.TXInputs
				for _, input := range inputs {
					if input.InputCanUnlock(address) {
						txo[string(input.TXID)] = append(txo[string(input.TXID)], input.Index)
					}
				}
			}

		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return utxos, total
}

//获取区块链
func GetBlockchian() *BlockChian {
	db, err := bolt.Open(blockchaindb, 0600, nil)
	if err != nil {
		log.Panic("open db err")

	}
	var tail []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if bucket == nil {
			log.Panic("create bucket err")
		}

		tail = bucket.Get([]byte("lastHash"))
		return nil
	})
	bc := BlockChian{
		db:db,
		tail:tail,
	}
	return &bc
}