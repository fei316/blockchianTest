package main

func main() {
	bc := NewBlockchian()

	cli := CLI{
		bc:bc,
	}

	cli.Run()





}
