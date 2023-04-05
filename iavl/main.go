package main

import (
	// "bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/bits"

	// "math/big"
	// "testing"
	types "github.com/tendermint/tendermint/types"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	leafPrefix  = []byte{0}
	innerPrefix = []byte{1}
)

// returns tmhash(0x00 || leaf)
func leafHash(leaf []byte) []byte {
	return tmhash.Sum(append(leafPrefix, leaf...))
}

// returns tmhash(0x01 || left || right)
func innerHash(left []byte, right []byte) []byte {
	data := make([]byte, len(innerPrefix)+len(left)+len(right))
	n := copy(data, innerPrefix)
	n += copy(data[n:], left)
	copy(data[n:], right)
	return tmhash.Sum(data)
}

// getSplitPoint returns the largest power of 2 less than length
func getSplitPoint(length int64) int64 {
	if length < 1 {
		panic("Trying to split a tree with size < 1")
	}
	uLength := uint(length)
	bitlen := bits.Len(uLength)
	k := int64(1 << uint(bitlen-1))
	if k == length {
		k >>= 1
	}
	return k
}

// synthetic txs to Txs
func b64ToHex(txs ...string) types.Txs {
	var txHexes []types.Tx
	for _, tx := range txs {
		// decode base64
		b, err := base64.StdEncoding.DecodeString(tx)
		if err != nil {
			panic(err)
		}
		txHexes = append(txHexes, b)
	}
	return txHexes
}

func txHashToBytes(txHashs []string) types.Txs {
	var txs types.Txs
	for _, txHash := range txHashs {
		txBytes, err := hex.DecodeString(txHash)
		if err != nil {
			panic(err)
		}
		txs = append(txs, txBytes)

	}
	return txs
}

func main() {

	txsHex := []string{
		"CocCCoICCiQvY29zbXdhc20ud2FzbS52MS5Nc2dFeGVjdXRlQ29udHJhY3QS2QEKK29yYWkxN3ZzZ3lmZmRnazUyNzlxeXV6dHZ6amd0MnE1a2x6aGpmNzIydzkSK29yYWkxOXA0M3kwdHFucjVxbGhmd254ZnQydTV1bnBoNXluNjB5N3R1dnUafXsid2l0aGRyYXciOnsiYXNzZXRfaW5mbyI6eyJuYXRpdmVfdG9rZW4iOnsiZGVub20iOiJpYmMvQTJFMkVFQzkwNTdBNEExQzJDMEE2QTRDNzhCMDIzOTExOERGNUYyNzg4MzBGNTBCNEE2QkREN0E2NjUwNkI3OCJ9fX19EgASZgpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA7tDviTE5kBZ7gxAHfh6xzZd9iyJs3YK+OUz4+L3oozlEgQKAggBGL4BEhEKCwoEb3JhaRIDNTAwEJbBEBpAV4sY7oyq2wgMNNHZfbrJwoAS71ectOC9BIPHZrf+kVQMHHD7ri8Pjz4ICQqiF+FWv8SdjnaQUJ2mmAA1o56QXg==",
		"CoICCv8BCiQvY29zbXdhc20ud2FzbS52MS5Nc2dFeGVjdXRlQ29udHJhY3QS1gEKK29yYWkxa21qcmxkZ2ozd2FrZjRxbWV1ZHJjZWQwbTl5N3FoMG01YXNmMngSK29yYWkxbmQ0cjA1M2Uza2dlZGdsZDJ5bWVuOGw5eXJ3OHhwanlhYWw3ajUaensiaW5jcmVhc2VfYWxsb3dhbmNlIjp7ImFtb3VudCI6Ijk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OTk5OSIsInNwZW5kZXIiOiJvcmFpMXlubWQyY2VtcnloY3d0anEzYWRoY3dheXJtODlsMmNyNHR3czR2In19EmUKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIS0r5ZO0japSdfqZC7mT4u85eoP1xyHn3n/x2lFnWHlxIECgIIARgCEhEKCwoEb3JhaRIDNzM2EO79CBpAJRZ6ytb+q1iTUkasCMly7iF+twZYIHEbUtUMpJfabKBkq5GKPPuaezY5k48ivQknLdyg5lu6ojv21NXamaPSVA==",
	}

	txs := b64ToHex(txsHex...)
	root := txs.Hash()

	fmt.Printf("%x\n", root)

	// txsBytes := [][]byte{
	// 	[]byte(txs[0]),
	// }

	txBzs := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		txBzs[i] = txs[i].Hash()
	}
	root2 := merkle.HashFromByteSlices(txBzs)
	fmt.Printf("%x\n", root2)

	root3, proofs := merkle.ProofsFromByteSlices(txBzs)
	fmt.Printf("%x\n", root3)
	// total
	fmt.Println("total :", proofs[0].Total)
	// index
	fmt.Println("index :", proofs[0].Index)
	// proof
	fmt.Println("proof :", proofs[0])
	// leaf hash
	fmt.Printf("leaf hash : %x\n", proofs[0].LeafHash)
	// aunts
	fmt.Printf("aunts : %x\n", proofs[0].Aunts[0])

	// // bytes of a charactor
	// fmt.Printf("%x\n", []byte("a"))
	// hashA := leafHash([]byte("a"))
	// fmt.Printf("hash of a: %x\n", hashA)

	// hashB := innerHash([]byte("a"), []byte("b"))
	// fmt.Printf("hash of b: %x\n", hashB)

	// sumHashAB := tmhash.Sum([]byte("a"))
	// fmt.Printf("hash of a: %x\n", sumHashAB)

}
