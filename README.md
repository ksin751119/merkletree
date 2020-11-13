<h1 align="center">Merkle Tree in Golang</h1>
<p align="center">
<a href="https://travis-ci.org/ksin751119/merkletree"><img src="https://travis-ci.org/ksin751119/merkletree.svg?branch=master" alt="Build"></a>
<a href="https://goreportcard.com/report/github.com/ksin751119/merkletree"><img src="https://goreportcard.com/badge/github.com/ksin751119/merkletree?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/ksin751119/merkletree"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

An implementation of a Merkle Tree written in Go. A Merkle Tree is a hash tree that provides an efficient way to verify
the contents of a set data are present and untampered with.

At its core, a Merkle Tree is a list of items representing the data that should be verified. Each of these items
is inserted into a leaf node and a tree of hashes is constructed bottom up using a hash of the nodes left and
right children's hashes. This means that the root node will effictively be a hash of all other nodes (hashes) in
the tree. This property allows the tree to be reproduced and thus verified by on the hash of the root node
of the tree. The benefit of the tree structure is verifying any single content entry in the tree will require only
nlog2(n) steps in the worst case.

#### Documentation

See the docs [here](https://godoc.org/github.com/ksin751119/merkletree).

#### Install
```
go get github.com/ksin751119/merkletree
```

#### Example Usage
Below is an example that makes use of the entire API - its quite small.
```go
package main

import (
  "crypto/sha256"
  "log"

  "github.com/ksin751119/merkletree"
)

//TestContent implements the Content interface provided by merkletree and represents the content stored in the tree.
type TestContent struct {
  x string
}

//CalculateHash hashes the values of a TestContent
func (t TestContent) CalculateHash() ([]byte, error) {
  h := sha256.New()
  if _, err := h.Write([]byte(t.x)); err != nil {
    return nil, err
  }

  return h.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (t TestContent) Equals(other merkletree.Content) (bool, error) {
  return t.x == other.(TestContent).x, nil
}

func main() {
  //Build list of Content to build tree
  var list []merkletree.Content
  list = append(list, TestContent{x: "Hello"})
  list = append(list, TestContent{x: "Hi"})
  list = append(list, TestContent{x: "Hey"})
  list = append(list, TestContent{x: "Hola"})

  //Create a new Merkle Tree from the list of Content
  t, err := merkletree.NewTree(list)
  if err != nil {
    log.Fatal(err)
  }

  //Get the Merkle Root of the tree
  mr := t.MerkleRoot()
  log.Println(mr)

  //Verify the entire tree (hashes for each node) is valid
  vt, err := t.VerifyTree()
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Verify Tree: ", vt)

  //Verify a specific content in in the tree
  vc, err := t.VerifyContent(list[0])
  if err != nil {
    log.Fatal(err)
  }

  log.Println("Verify Content: ", vc)

  //String representation
  log.Println(t)
}

```

#### Example Usage 2
Below is an example that makes use of the entire API - its quite small.

```go
package main

import (
  "crypto/sha256"
  "log"
  "golang.org/x/crypto/sha3"

  "github.com/ksin751119/merkletree"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/crypto"

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

func hashSort(left int, leftHash []byte, right int, rightHash[]byte) (int, int){
	if bytes.Compare(nl[left].Hash, nl[right].Hash) > 0 {
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

  config := &TreeConfig{
    	HashStrategy: sha3.NewLegacyKeccak256,
	    HashSortFunc: hashSort,
  }

	//Create a new Merkle Tree from the list of Content
	t, err := merkletree.NewTreeWithConfig(list, sha3.NewLegacyKeccak256)

	if err != nil {
		log.Fatal(err)
	}

	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()
	log.Printf("Root: 0x%x", mr)

	proofs, deep, _ := t.GetMerklePath(list[1])
	for idx, pr1 := range proofs {
		log.Printf("proof[%d]: 0x%x", idx, pr1)
	}

	for idx, l := range list {
		h, _ := l.CalculateHash()
		log.Printf("leaf[%d] hash: %x\n", idx, h)
	}
}

```


#### License
This project is licensed under the MIT License.
