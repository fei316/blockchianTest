package main

import (
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type BlockchainIterator struct {
	db                 *bolt.DB
	currentPointerHash []byte
}

func (blockchain *BlockChian) NewBlockchainInterator() (*BlockchainIterator) {


	return &BlockchainIterator{
		db:                 blockchain.db,
		currentPointerHash: blockchain.Tail,
	}


}

func (bcIterator *BlockchainIterator) Next() *Block {
	var block Block
	db, err := bolt.Open(blockchaindb, 0600, &bolt.Options{Timeout: 5 * time.Second})
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockbucket))
		if bucket == nil {
			log.Panic("bucket为空，不应该为空")
		}

		var temphash = bcIterator.currentPointerHash


		blocktmp := bucket.Get(temphash)
		BUGAGAIN:
		for {
			if len(blocktmp) == 0 {
				err = db.View(func(tx *bolt.Tx) error {
					bucket := tx.Bucket([]byte(blockbucket))
					if bucket == nil {
						log.Panic("create bucket err")
					}

					tail := bucket.Get([]byte("lastHash"))
					temphash = tail
					return nil
				})
				if err != nil {
					log.Panic(err)
				}
				blocktmp = bucket.Get(temphash)

				continue BUGAGAIN
			} else {
				break
			}
		}

		block = DeSerialize(blocktmp)


		bcIterator.currentPointerHash = block.PrevHash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &block
}
