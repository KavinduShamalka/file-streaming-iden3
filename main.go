package main

import (
	"context"
	"encoding/json"
	"fmt"
	"iden3-test/streaming"
	"math/big"
	"os"
	"strconv"

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
	chunk, hasher, err := streaming.SplitFile(inputFile, chunkSize)
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

	// add chunks to merkle tree
	for index, value := range chunk {

		// Convert string to int64
		int64Value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			// Handle error if conversion fails
			fmt.Println("Conversion error:", err)
			return
		}

		mt.Add(ctx, big.NewInt(int64(index)), big.NewInt(int64Value))

		// Proof of membership of a leaf with index
		proofExist, _, _ := mt.GenerateProof(ctx, big.NewInt(int64(index)), mt.Root())
		fmt.Printf("Proof of membership %v: %v\n", index, proofExist.Existence)

		err = newFunction(proofExist, chunk, hasher, "restored_data.jpg")
		if err != nil {
			fmt.Println("Error retrieving and verifying chunks:", err)
			return
		}

	}

	fmt.Printf("Root tree address: %v\n", mt)

	// transform root from bytes array to json
	root, _ := json.Marshal(mt.Root().BigInt())

	fmt.Println(string(root))

}

func newFunction(proofExist *merkletree.Proof, chunkNames []string, hashValues []string, outputFileName string) error {
	if proofExist.Existence {
		return streaming.RetrieveChunksAndVerify(chunkNames, hashValues, outputFileName)
	}
	return fmt.Errorf("proof of non-membership received")
}
