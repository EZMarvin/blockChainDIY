package cli

import (
	"flag"
	"fmt"
	"github.com/EZMarvin/go-blockchain/blockchain"
	"github.com/EZMarvin/go-blockchain/wallet"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{
	//blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get balance for that address")
	fmt.Println(" createblockchain -address ")
	fmt.Println(" printchain - Prints the blocks in chain")
	fmt.Println(" send - from From -to To -amount AMOUNT - Send amount")
	//fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	//fmt.Println(" print - Prints the blocks in the chain")
	fmt.Println("createwallet - Create a new Wallet")
	fmt.Println("listaddresses = Lists the addresses in out wallet file")
}

func (cli *CommandLine) validateArgs(){
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

/*func (cli *CommandLine) addBlock(data string){
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block")
}*/

func (cli *CommandLine) printChain(){
	chain := blockchain.ContinueBlocChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		//fmt.Printf("Data: %s\n", block.Data)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlocChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance )
}

func (cli *CommandLine) send(from, to string, amount int){
	chain := blockchain.ContinueBlocChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success send")
}
func (cli *CommandLine) createWallet(){
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is %s\n", address)
}

func (cli *CommandLine) listAddresses(){
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listaddressCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")


	/*
		addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
		printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
		addBlockData := addBlockCmd.String("block", "", "Block data")
	*/
	switch os.Args[1] {
	/*
		case "add":
			err := addBlockCmd.Parse(os.Args[2:])
			blockchain.Handle(err)

		case "print":
			err := printChainCmd.Parse(os.Args[2:])
			blockchain.Handle(err)
	*/
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listaddressCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listaddressCmd.Parsed() {
		cli.listAddresses()
	}
	/*
		if addBlockCmd.Parsed() {
			if *addBlockData == "" {
				addBlockCmd.Usage()
				runtime.Goexit()
			}
			cli.addBlock(*addBlockData)
		}

		if printChainCmd.Parsed() {
			cli.printChain()
		}
	*/

}