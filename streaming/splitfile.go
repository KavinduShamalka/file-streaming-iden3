package streaming

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// SplitFile splits a file into chunks of a specified size
func SplitFile(inputFile *os.File, chunkSize int64) ([]string, []string, error) {

	// // Open the input file
	// file, err := os.Open(inputFile)

	// if err != nil {
	// 	return nil, err
	// }

	// Close the file, when the function exits
	defer inputFile.Close()

	// Get file info
	fileInfo, err := inputFile.Stat()
	if err != nil {
		return nil, nil, err
	}

	fileSize := fileInfo.Size()

	chunkNames := make([]string, 0)
	// chunkSizes := make([]int64, 0)
	hashValues := make([]string, 0)

	hasher := sha256.New()

	for i := int64(0); i < fileSize; i += chunkSize {
		chunkName := fmt.Sprintf("%v_chunk%d", inputFile, i/chunkSize+1)

		chunkFile, err := os.Create(chunkName)
		if err != nil {
			return nil, nil, err
		}

		// Create a multi-writer to both write to the file and calculate the hash
		multiWriter := io.MultiWriter(chunkFile, hasher)

		// Copy the chunkSize bytes from the original file to the chunk file
		_, err = io.CopyN(multiWriter, inputFile, chunkSize)
		if err != nil && err != io.EOF {
			return nil, nil, err
		}

		chunkFile.Close()
		chunkNames = append(chunkNames, chunkName)
		// chunkSizes = append(chunkSizes, chunkSize)

		// Add the hash value of the chunk to the slice
		hashValue := fmt.Sprintf("%x", hasher.Sum(nil))

		hashValues = append(hashValues, hashValue)
		newChunkName := hashValue
		err = os.Rename(chunkName, newChunkName)
		if err != nil {
			return nil, nil, err
		}

		// chunkNames = append(chunkNames, newChunkName)
		hashValues = append(hashValues, hashValue)
		// Reset the hash for the next iteration
		hasher.Reset()

	}

	return chunkNames, hashValues, nil
}
