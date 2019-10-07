package main

import (
	"fmt"
	"log"
)


func (cli *CLI) printChain() {
	bc := cli.bc

	bcIterator:= bc.NewBlockchainInterator()
	fmt.Println("*************区块链遍历开始*************")
	for {
		block := bcIterator.Next()
		fmt.Printf("=================================================================================\n")
		fmt.Printf("版本号：%d\n", block.Version)
		fmt.Printf("前Hash：%x\n", block.PrevHash)
		fmt.Printf("梅克尔根：%x\n", block.MerkelRoot)
		fmt.Printf("时间戳：%d\n", block.TimeStamp)
		fmt.Printf("难度值：%d\n", block.Difficulty)
		fmt.Printf("随机数：%d\n", block.Nonce)
		fmt.Printf("当前Hash：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", string(block.Transactions[0].TXInputs[0].Signature[:]))

		if len(block.PrevHash) == 0 {

			break
		}
	}
	fmt.Println("*************区块链遍历结束*************")
}

func (cli *CLI) getBalance(address string) {

	if !IsValidAddress(address) {
		log.Printf("地址无效[%s]\n", address)
		return
	}
	pubHash := GetPubHashByAddress(address)

	utxos := cli.bc.getUTXOs(pubHash)

	var total = 0.0
	for _, utxo := range utxos {

		total = total + utxo.Value
	}
	fmt.Printf("地址：[%s]的余额为：%f\n", address, total)
}

//交易
func (cli *CLI) send(from, to string, amount float64, miner string, remark string)  {

	if !IsValidAddress(from) {
		log.Panic("地址%s无效", from)
		return
	}
	if !IsValidAddress(to) {
		log.Panic("地址%s无效", to)
		return
	}
	if !IsValidAddress(miner) {
		log.Panic("地址%s无效", miner)
		return
	}

	bc := GetBlockchian()
	var trans []*Transaction
	coinbase := NewCoinbaseTx(miner, remark)

	trans = append(trans, coinbase)

	tran := NewTransaction(from, to, amount, bc)

	trans = append(trans, tran)

	bc.AddBlock(trans)
}

//创建钱包
func (cli *CLI) createWalet() {
	ws := NewWallets()
	address := ws.createWallet()
	fmt.Printf("地址创建成功：%s\n", address)
}

//列出钱包所有地址
func (cli *CLI)listAddrs()  {
	ws := NewWallets()
	fmt.Printf("******地址开始******\n")
	for address, _ := range ws.WalletsMap {
		fmt.Printf("[%s]\n",string(address))
	}
	fmt.Printf("******地址结束******\n")
}