package main


//区块链结构体
type BlockChian struct {
	blocks []Block
}

//创建区块链
func NewBlockchian() *BlockChian {
	blockchian := BlockChian{
		blocks:[]Block{*GenesisBlock()},
	}
	return &blockchian
}

//区块链添加区块
func (blockchain *BlockChian) AddBlock(data string) {
	lastBlock := blockchain.blocks[len(blockchain.blocks) - 1]
	block := NewBloack(data, lastBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, *block)
}
