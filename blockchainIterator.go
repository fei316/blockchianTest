package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	db                 *bolt.DB
	currentPointerHash []byte
}

func (blockchain *BlockChian) NewBlockchainInterator() *BlockchainIterator {

	return &BlockchainIterator{
		db:                 blockchain.db,
		currentPointerHash: blockchain.tail,
	}

}

func (bcIterator *BlockchainIterator) Next() *Block {
	var block Block
	db := bcIterator.db
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if bucket == nil {
			log.Panic("bucket为空，不应该为空")
		}
		blocktmp := bucket.Get(bcIterator.currentPointerHash)

		block = DeSerialize(blocktmp)
		bcIterator.currentPointerHash = block.PrevHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &block
}
