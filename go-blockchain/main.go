package main

import (
	"github.com/EZMarvin/go-blockchain/cli"
	"os"

)



func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	/*
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	*/
	cli.Run()
}


