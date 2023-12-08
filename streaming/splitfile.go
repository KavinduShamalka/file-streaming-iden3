package streaming

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// SplitFile splits a file into chunks of a specified size
func SplitFile(inputFile *os.File, chunkSize int64) ([]string, []string, error) {

	// close the input file after the end of the split file func execute
	defer inputFile.Close()

	// Stat returns the FileInfo structure describing file. If there is an error, it will be of type *PathError.
	fileInfo, err := inputFile.Stat()
	if err != nil {
		return nil, nil, err
	}

	// Get the file size
	fileSize := fileInfo.Size()
	var chunkNames []string //create chunk names slice
	var hashValues []string //create hash values slice

	// Initate new hasher
	hasher := sha256.New()

	// For loop for split file into the chunks
	for i := int64(0); i < fileSize; i += chunkSize {

		// create file chunks
		chunkFile, err := os.Create(fmt.Sprintf("%v_chunk%d", inputFile.Name(), i/chunkSize+1))
		if err != nil {
			return nil, nil, err
		}

		// creates a writer that duplicates its writes to all the provided writers
		multiWriter := io.MultiWriter(chunkFile, hasher)

		_, err = io.CopyN(multiWriter, inputFile, chunkSize)
		if err != nil && err != io.EOF {
			return nil, nil, err
		}

		chunkFile.Close()

		// Calculate hash value
		hashValue := fmt.Sprintf("%x", hasher.Sum(nil))
		hashedFileName := hashValue

		// Rename the chunk file with its hash value
		err = os.Rename(chunkFile.Name(), hashedFileName)
		if err != nil {
			return nil, nil, err
		}

		chunkNames = append(chunkNames, hashedFileName)
		hashValues = append(hashValues, hashValue)

		hasher.Reset()
	}

	return chunkNames, hashValues, nil
}
