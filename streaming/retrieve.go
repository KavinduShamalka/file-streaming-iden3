package streaming

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Retrieve Chunks and Verify.
func RetrieveChunksAndVerify(chunkNames []string, outputFileName string) error {

	// Create output file name
	outputFile, err := os.Create(outputFileName)

	// Check if there is any error
	if err != nil {
		return err
	}

	// Defer is used to ensure that a function call is performed later in a program’s execution, usually for purposes of cleanup.
	// defer: Allows us to "clean up" resources even if the function encounters errors.
	// close: Can return errors indicating problems with resource release, which can be handled appropriately.
	defer outputFile.Close()

	// Use sha256 as the hashing algorithum
	hasher := sha256.New()

	// get the index and the chunknames from "chunkNames" slice
	for i, chunkName := range chunkNames {
		chunkFile, err := os.Open(chunkName)
		if err != nil {
			return err
		}

		defer chunkFile.Close()

		// Create a multi-reader to both read from the file and calculate the hash
		multiReader := io.TeeReader(chunkFile, hasher)

		_, err = io.Copy(outputFile, multiReader)
		if err != nil {
			return err
		}

		// Verify the hash of the chunk
		hashValue := fmt.Sprintf("%x", hasher.Sum(nil))

		// alter file name
		// hashValues[2] = "sithum"

		//check both hash values are same
		if hashValue != chunkNames[i] {
			return fmt.Errorf("hash verification failed for chunk %s", chunkName)
		}

		// Reset the hash for the next iteration
		hasher.Reset()
	}

	return nil
}
