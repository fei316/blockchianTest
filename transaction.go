package main

type Transaction struct {
	TXID []byte
	TXInputs []TXInput
	TXOutputs []TXOutput

}

type TXInput struct {
	TXID []byte
	Index int64
	sig string
}

type TXOutput struct {
	value float64
	PubkeyHash string
}