package main

import "fmt"

func (cli *CLI) addBlock(txs []*Transaction) {
	//TODO
	//cli.bc.AddBlock(data)
}

func (cli *CLI) printChain() {
	bc := cli.bc

	bcIterator := bc.NewBlockchainInterator()
	fmt.Println("*************区块链遍历开始*************")
	for {
		block := bcIterator.Next()
		fmt.Printf("=====================\n")
		fmt.Printf("版本号：%d\n", block.Version)
		fmt.Printf("前Hash：%x\n", block.PrevHash)
		fmt.Printf("梅克尔根：%x\n", block.MerkelRoot)
		fmt.Printf("时间戳：%d\n", block.TimeStamp)
		fmt.Printf("难度值：%d\n", block.Difficulty)
		fmt.Printf("随机数：%d\n", block.Nonce)
		fmt.Printf("当前Hash：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", block.Transactions[0].TXInputs[0].Sig)

		if len(block.PrevHash) == 0 {

			break
		}
	}
	fmt.Println("*************区块链遍历结束*************")
}

func (cli *CLI) getBalance(address string) {
	utxos := cli.bc.getUTXOs(address)
	var total = 0.0
	for _, utxo := range utxos {
		total = total + utxo.value
	}
	fmt.Printf("地址：[%s]的余额为：%f\n", address, total)
}
