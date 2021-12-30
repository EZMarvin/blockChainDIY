package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID []byte
	Inputs []TxInput
	Output []TxOutput
}

// SetID set Id for transaction
func (tx *Transaction) SetID(){
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction{
	if data == ""{
		data = fmt.Sprintf("Coin to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Error: not enough account")
	}

	// iterate through the valid outs and create tx's input
	for txid, outs := range validOutputs{
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _,out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc-amount, from})
	}

	tx := Transaction{nil,inputs,outputs}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *Transaction) Sign (privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinbase(){
		return
	}

	for _, in := range tx.Inputs {
		if prevTxs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Error: Previous transaction does not exist")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		pre
	}
}
