package main

import (
	"fmt"
	"os"
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
	default:
		fmt.Println(Usage)

	}
}


