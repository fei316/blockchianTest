package main

import (
	"log"
	"testing"
)

func TestNewProofOfWork(t *testing.T) {
	ts := Transaction{
		TXID:[]byte{},
	}
	log.Print(ts)

	block := NewBloack([]*Transaction{&ts}, []byte{})
	NewProofOfWork(block)




}
