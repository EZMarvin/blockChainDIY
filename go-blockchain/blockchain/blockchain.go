package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

const (
	dbPath = "./tmp/blocks"
	dbFile = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
	//Blocks []*Block
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database *badger.DB
}

func DBexist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err){
		return false
	}
	return true
}
func InitBlockChain(address string) *BlockChain{
	var lastHash []byte

	if DBexist() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error{
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash
		return err
	})

	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlocChain(address string) *BlockChain {
	if DBexist() == false {
		fmt.Println("no existing blockchain found, use init to creat one")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error{
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})

	Handle(err)
	chain := BlockChain{lastHash, db}
	return &chain
}

func (chain *BlockChain) AddBlock(transactions []*Transaction){
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})

	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error{
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	Handle(err)
}

func (blockchain *BlockChain) Iterator() *BlockChainIterator{
	return &BlockChainIterator{blockchain.LastHash, blockchain.Database}
}

func (iter *BlockChainIterator) Next() *Block{
	var block *Block

	err := iter.Database.View(func (txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodeBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodeBlock)
		return err
	})

	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _,tx := range block.Transactions{
			txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx, out := range tx.Output {
					if spentTXOs[txID] != nil {
						for _, spentOut := range spentTXOs[txID] {
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}
					if out.CanBeUnlocked(address) {
						unspentTxs = append(unspentTxs, *tx)
					}
				}
				if tx.IsCoinbase() == false {
					for _, in := range tx.Inputs {
						if in.CanUnlock(address){
							inTxID := hex.EncodeToString(in.ID)
							spentTXOs[inTxID] = append(spentTXOs[inTxID],in.Out)
						}
					}
				}
		}
		if len(block.PrevHash) == 0 {
			break
		}

	}
	return unspentTxs
}

// FindUTXO get total balance of user account
func (chain *BlockChain) FindUTXO (address string) []TxOutput {
	var UTXO []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Output {
			if out.CanBeUnlocked(address) {
				UTXO = append(UTXO, out)
			}
		}
	}

	return UTXO
}

// FindSpendableOutputs check user have enough coin balance for creating transaction

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulate := 0

	Work:
		for _,tx := range unspentTxs{
			txID := hex.EncodeToString(tx.ID)
			for outIdx, out := range tx.Output{
				accumulate += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulate >= amount{
					break Work
				}
			}
		}
	return accumulate, unspentOuts
}
