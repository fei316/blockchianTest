package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChian
}

const Usage = `
	print  "打印区块链"
	balance --address Address "查询地址余额"
	send --from FROM --to TO --amount AMOUNT --miner MINER DATA
	createWallet "创建钱包"
	listAddrs "查看地址"
`

func (cli *CLI) Run() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}

	switch args[1] {
	case "print":
		cli.printChain()
		break
	case "balance":
		if len(args) == 4 && args[2] == "--address" {
			cli.getBalance(args[3])
		} else {
			fmt.Printf(Usage)
		}
		break
	case "send":	//send --from FROM --to TO --amount AMOUNT --miner MINER "标识"
		if len(args) == 11 {
			amount, err := strconv.ParseFloat(args[7],64)
			if err != nil {
				fmt.Printf(Usage)
			} else {
				cli.send(args[3], args[5], amount, args[9], args[10])
			}

		} else {
			fmt.Printf(Usage)
		}
		break
	case "createWallet":	//createWallet
		cli.createWalet()
		break
	case "listAddrs":	//listAddrs

		cli.listAddrs()
		break
	default:
		fmt.Println(Usage)

	}
}


