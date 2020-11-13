package main

import (
	"bytes"
	"log"
	"math/big"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ksin751119/merkletree"
)

type TestContent struct {
	Address string
	Amount  *big.Int
}

func (t TestContent) CalculateHash() ([]byte, error) {
	hash := crypto.Keccak256(
		common.HexToAddress(t.Address).Bytes(),
		common.LeftPadBytes(t.Amount.Bytes(), 32),
	)
	return hash, nil
}

func (t TestContent) Equals(other merkletree.Content) (bool, error) {
	return t.Address == other.(TestContent).Address && t.Amount == other.(TestContent).Amount, nil

}

func hashSort(left int, leftHash []byte, right int, rightHash []byte) (int, int) {
	if bytes.Compare(leftHash, rightHash) > 0 {
		return right, left
	}
	return left, right
}

func main() {
	//Build list of Content to build tree
	var list []merkletree.Content
	list = append(list, TestContent{Address: "0x30afBFe6B5eBC2F5f008F819fc0Eb1E71ad5B265", Amount: big.NewInt(1000000000000000000)})
	list = append(list, TestContent{Address: "0x1b57b3A1d5b4aa8E218F54FafB00975699463e6e", Amount: big.NewInt(1000000000000000000)})
	list = append(list, TestContent{Address: "0xAA293A146aAf9E05BeDD1Ff29B0da5bD8BE70955", Amount: big.NewInt(1000000000000000000)})

	config := &merkletree.TreeConfig{
		HashStrategy: sha3.NewLegacyKeccak256,
		HashSortFunc: hashSort,
	}

	//Create a new Merkle Tree from the list of Content
	t, err := merkletree.NewTreeWithConfig(list, config)

	if err != nil {
		log.Fatal(err)
	}

	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()
	log.Printf("Root: 0x%x", mr)

	proofs, _, _ := t.GetMerklePath(list[1])
	for idx, pr1 := range proofs {
		log.Printf("proof[%d]: 0x%x", idx, pr1)
	}

	for idx, l := range list {
		h, _ := l.CalculateHash()
		log.Printf("leaf[%d] hash: %x\n", idx, h)
	}
}
