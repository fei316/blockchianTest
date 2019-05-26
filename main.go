package main

func main() {
	bc := NewBlockchian()
	bc.AddBlock("小明向向小红转了100个比特币")
	/*for i, block := range bc.blocks {
		fmt.Printf("=====当前区块高度：%d =======\n", i)
		fmt.Printf("前Hash：%x\n", block.PrevHash)
		fmt.Printf("当前Hash：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", block.Data)
	}*/

}
