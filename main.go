package main

import (
	"context"
	"encoding/json"
	"fmt"
	"iden3-test/streaming"
	"math/big"
	"os"

	merkletree "github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
)

func main() {

	inputFile, err := os.Open("test.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	chunkSize := int64(256000) // Set your desired chunk size in bytes

	// Split the file into chunks
	chunk, err := streaming.SplitFile(inputFile, chunkSize)
	if err != nil {
		fmt.Println("Error splitting file:", err)
		return
	}

	// Sparse MT
	ctx := context.Background()

	// Tree storage
	store := memory.NewMemoryStorage()

	// Generate a new MerkleTree with 32 levels
	mt, _ := merkletree.NewMerkleTree(ctx, store, 32)

	// Add a leaf to the tree with index 1 and chunk[0]
	// index1 := big.NewInt(1)
	// value1 := big.NewInt(chunk[0])
	// mt.Add(ctx, index1, value1)

	// // Add another leaf to the tree
	// index2 := big.NewInt(2)
	// value2 := big.NewInt(chunk[1])
	// mt.Add(ctx, index2, value2)

	for i, val := range chunk {
		mt.Add(ctx, big.NewInt(int64(i)), big.NewInt(val))
	}

	// Proof of membership of a leaf with index 1
	proofExist, value, _ := mt.GenerateProof(ctx, big.NewInt(0), mt.Root())

	fmt.Println("Proof of membership 1:", proofExist.Existence)
	fmt.Println("Value corresponding to the queried index:", value)

	// Proof of non-membership of a leaf with index 4
	proofNotExist, _, _ := mt.GenerateProof(ctx, big.NewInt(6), mt.Root())

	fmt.Println("Proof of membership 4:", proofNotExist.Existence)

	fmt.Printf("%v\n", mt)

	// transform root from bytes array to json
	claimToMarshal, _ := json.Marshal(mt.Root())

	fmt.Println(string(claimToMarshal))

}
