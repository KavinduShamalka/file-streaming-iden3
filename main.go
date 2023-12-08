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

	// Define chuksize 256kb in bytes
	chunkSize := int64(256000)

	// Read the input file
	inputFile, err := os.Open("test.png")

	// Check if there is an error
	if err != nil {
		// Print the error
		fmt.Printf("Error: %v", err)
		return
	}

	// Used to delay the execution of a function until the surrounding function completes.
	defer inputFile.Close()

	// Return the chukn name, and hash value from the SplitFile function.
	chunkNames, hashValues, err := streaming.SplitFile(inputFile, chunkSize)
	if err != nil {
		fmt.Println("Error splitting and hashing file:", err)
		return
	}

	// Create a Context Background
	ctx := context.Background()

	// Declare new memory
	store := memory.NewMemoryStorage()

	// Create merkle tree
	mt, _ := merkletree.NewMerkleTree(ctx, store, 32)

	// Get index and value of from the ChunkNames slice
	for index, value := range chunkNames {

		//Add to the merkle tree
		mt.Add(ctx, big.NewInt(int64(index)), big.NewInt(0))
		fmt.Println()
		fmt.Println(ctx, index, value)

		fmt.Println(mt.Root())

		// Proof of membership for each chunk
		proofExist, _, _ := mt.GenerateProof(ctx, big.NewInt(int64(index)), mt.Root())
		fmt.Printf("Proof of membership for chunk %d: %v\n", index, proofExist.Existence)

		//Check the proof
		err := checkProof(proofExist, chunkNames, hashValues, "restored_data.jpg")
		if err != nil {
			fmt.Println("Error retrieving and verifying chunks:", err)
			return
		}
	}

	// Proof of non-membership for a non-existing chunk (e.g., index 100)
	nonExistingIndex := big.NewInt(100)                                       // Intialize none existing index
	proofNotExist, _, _ := mt.GenerateProof(ctx, nonExistingIndex, mt.Root()) //Generate the proof
	fmt.Printf("Proof of non-membership for chunk %d: %v\n", nonExistingIndex.Int64(), proofNotExist.Existence)

	claimToMarshal, _ := json.Marshal(mt.Root())
	fmt.Println(string(claimToMarshal))

}

// Check the chunk's proof of existence
func checkProof(proofExist *merkletree.Proof, chunkNames []string, hashValues []string, outputFileName string) error {

	// check if proofExist true
	if proofExist.Existence {
		// RetrieveChunksfunctions
		return streaming.RetrieveChunksAndVerify(chunkNames, hashValues, outputFileName)
	}
	return fmt.Errorf("proof of non-membership received")
}
