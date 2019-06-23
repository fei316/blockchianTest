package main

func main() {
	bc := NewBlockchian("北京市海淀区安宁华庭3区5号楼2单元702")

	cli := CLI{
		bc: bc,
	}

	cli.Run()

}
