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
	add --data DATA "添加区块到区块链"
	print 		    "打印区块链"
	balance --address Address "查询地址余额"
`

func (cli *CLI) Run() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}

	switch args[1] {
	case "add":
		if len(args) == 4 && args[2] == "--data" {
			//TODO
			//cli.addBlock(args[3])
		} else {
			fmt.Printf(Usage)
		}
		break
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
	case "send":	//send --from FROM --to TO --amount AMOUNT --miner MINER
		if len(args) == 10 {
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

		break
	default:
		fmt.Println(Usage)

	}
}


