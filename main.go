package main

import "fmt"

func main() {
	bc := NewBlockchian()
	bc.AddBlock("小明向向小红转了100个比特币")
	bc.AddBlock("小明向向小红转了200个比特币")
	bc.AddBlock("小红向小刚转了200个比特币")

	bcIterator := bc.NewBlockchainInterator()
	fmt.Println("*************区块链遍历开始*************")
	for {
		block := bcIterator.Next()
		fmt.Printf("=====================\n")
		fmt.Printf("前Hash：%x\n", block.PrevHash)
		fmt.Printf("当前Hash：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", block.Data)

		if (len(block.PrevHash) == 0) {

			break;
		}
	}
	fmt.Println("*************区块链遍历结束*************")




}
