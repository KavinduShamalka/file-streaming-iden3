package streaming

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// SplitFile splits a file into chunks of a specified size
func SplitFile(inputFile *os.File, chunkSize int64) ([]string, []string, error) {

	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return nil, nil, err
	}

	fileSize := fileInfo.Size()
	var chunkNames []string
	var hashValues []string

	hasher := sha256.New()

	for i := int64(0); i < fileSize; i += chunkSize {
		chunkFile, err := os.Create(fmt.Sprintf("%v_chunk%d", inputFile.Name(), i/chunkSize+1))
		if err != nil {
			return nil, nil, err
		}

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
