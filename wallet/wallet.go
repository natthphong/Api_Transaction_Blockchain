package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	///create key
	w := new(Wallet)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {

		panic(err)
	}
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	//hash public  key
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// ripemd160 hash
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	//add version byte
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	//
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	checksum := digest6[:4]
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4)
	copy(dc8[21:], checksum[:])

	address := base58.Encode(dc8)
	w.blockchainAddress = address
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}
func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

type Transaction struct {
	senderPrivate            *ecdsa.PrivateKey
	senderPublic             *ecdsa.PublicKey
	senderBlockchainAddress  string
	recipienBlockchainAdress string
	value                    float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string,
	recipien string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipien, value}
}

func (t *Transaction) GenerateSignature() *Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, err := ecdsa.Sign(rand.Reader, t.senderPrivate, h[:])

	if err != nil {
		log.Print(err.Error())
	}

	return &Signature{r, s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Sender    string  `json:sender_blockchain_address`
			Recipient string  `json:recipient_blockchain_address`
			Value     float32 `value`
		}{
			Sender:    t.senderBlockchainAddress,
			Recipient: t.recipienBlockchainAdress,
			Value:     t.value,
		},
	)
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
