package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	db *bolt.DB
	currentPointerHash []byte
}

func (bc *BlockChian) NewBlockchainInterator() *BlockchainIterator {

	return &BlockchainIterator{
		db:bc.db,
		currentPointerHash:bc.tail,
	}

}

func (bcIterator *BlockchainIterator) Next() *Block {
	var block Block
	db := bcIterator.db
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if (bucket == nil) {
			log.Panic("bucket为空，不应该为空")
		}
		blocktmp := bucket.Get(bcIterator.currentPointerHash)

		block = DeSerialize(blocktmp)
		bcIterator.currentPointerHash = block.PrevHash
		return nil
	})
	return &block
}
