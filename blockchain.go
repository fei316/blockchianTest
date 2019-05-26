package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const blockchaindb = "blockchain.db"
const blockbucket = "blockbucket"


//区块链结构体
type BlockChian struct {
	db *bolt.DB
	tail []byte
}

//创建区块链
func NewBlockchian() *BlockChian  {

	db, err := bolt.Open(blockchaindb, 0600, nil)
	defer db.Close()
	if (err != nil) {
		log.Panic("open db err")
	}

	var tail = []byte{}

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if (bucket == nil) {
			bucket, err = tx.CreateBucket([]byte(blockbucket))
			if (err != nil) {
				log.Panic("create bucket err")
			}
			block := GenesisBlock()
			bucket.Put(block.Hash, []byte{})//TODO value需要写函数转换
			bucket.Put([]byte("lastHash"), block.Hash)
			tail = block.Hash
		}

		tail = bucket.Get([]byte(blockbucket))
		return nil
	})
	return &BlockChian{
		db:db,
		tail:tail,
	}
}

//区块链添加区块
func (blockchain *BlockChian) AddBlock(data string) {
	/*lastBlock := blockchain.blocks[len(blockchain.blocks) - 1]
	block := NewBloack(data, lastBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, *block)*/
}
