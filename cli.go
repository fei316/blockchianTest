package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChian
}

const Usage  =  `
	addBlock --data DATA "add block to blockchain"
	printChain 		     "print blockchain"
`

func (cli *CLI) Run() {
	args := os.Args

	if (len(args) < 2) {
		fmt.Println(Usage)
		return;
	}

	switch args[1] {
	case "addBlock":
		if (len(args) == 4 && args[2] == "--data") {
			cli.addBlock(args[3])
		} else {
			fmt.Printf(Usage)
		}
		break
	case "printChain":
		cli.printChain()
		break
	default:
		fmt.Println(Usage)

	}
}
