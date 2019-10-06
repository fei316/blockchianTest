package main

func main() {
	bc := NewBlockchian("testaddress")

	cli := CLI{
		bc: bc,
	}

	cli.Run()

}
