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
