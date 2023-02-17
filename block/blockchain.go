package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFDICULTY = 3
	MINIG_SENDER      = "THE BLOCKCHAIN"
	MINIG_REWARD      = 1.0
)

// Block
type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

// /func Block
func NewBlock(nonce int, prevaiusHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = prevaiusHash
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp           	%d\n", b.timestamp)
	fmt.Printf("nonce 					%d\n", b.nonce)
	fmt.Printf("previousHash            %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}

}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	//fmt.Print(string(m))
	return sha256.Sum256([]byte(m))
}
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{Timestamp: b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

// /BlockChain
type Blockchain struct {
	transactionPool  []*Transaction
	chain            []*Block
	blockchainAdress string
}

// Initchain
func NewBlockChain(blockchainAdress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAdress = blockchainAdress
	bc.CreateBlock(0, b.Hash()) //create block 0
	return bc
}

//function Blockchain

func (bc *Blockchain) CreateBlock(nonce int, previusHash [32]byte) *Block {
	b := NewBlock(nonce, previusHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s  Chain %d  %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 50))
}

// addtransaction

func (bc *Blockchain) AddTransaction(sender string, recipien string, value float32) {
	t := NewTransaction(sender, recipien, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)

	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderBlockchainAdress, t.recipienBlockchianAddress, t.value))

	}
	return transactions

}

// /ProofofWork
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulity int) bool {
	zeros := strings.Repeat("0", difficulity)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulity] == zeros
}

func (bc *Blockchain) ProofofWork() int {

	transaction := bc.CopyTransactionPool()

	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transaction, MINING_DIFDICULTY) {
		nonce += 1
	}
	return nonce
}

// /Minig generate Block
func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINIG_SENDER, bc.blockchainAdress, MINIG_REWARD)
	nonce := bc.ProofofWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining , status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipienBlockchianAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAdress {
				totalAmount -= value
			}

		}
	}
	return totalAmount
}

// struct and function Transaction
type Transaction struct {
	senderBlockchainAdress    string
	recipienBlockchianAddress string
	value                     float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	fmt.Printf("sender_blockchain_address		%s\n", t.senderBlockchainAdress)
	fmt.Printf("recipient_blockchain_address		%s\n", t.recipienBlockchianAddress)
	fmt.Printf("value					%.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAdress    string  `json:"senderBlockchainAdress"`
		RecipienBlockchianAddress string  `json:"recipienBlockchianAddress"`
		Value                     float32 `json:"value"`
	}{
		SenderBlockchainAdress:    t.senderBlockchainAdress,
		RecipienBlockchianAddress: t.recipienBlockchianAddress,
		Value:                     t.value,
	})
}
