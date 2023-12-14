package streaming

import (
	"fmt"
	"io"
	"os"
)

// SplitFile splits a file into chunks of a specified size
func SplitFile(inputFile *os.File, chunkSize int64) ([]string, error) {

	// close the input file after the end of the split file func execute
	defer inputFile.Close()

	// Stat returns the FileInfo structure describing file. If there is an error, it will be of type *PathError.
	fileInfo, err := inputFile.Stat()
	if err != nil {
		return nil, err
	}

	// Get the file size
	fileSize := fileInfo.Size()
	var chunkNames []string //create chunk names slice
	// var hashValues []string //create hash values slice

	// Initate new hasher
	// hasher := sha256.New()

	// For loop for split file into the chunks
	for i := int64(0); i < fileSize; i += chunkSize {

		// create file chunks
		chunkFile, err := os.Create(fmt.Sprintf("%v_chunk%d", inputFile.Name(), i/chunkSize+1))
		if err != nil {
			return nil, err
		}

		// creates a writer that duplicates its writes to all the provided writers
		multiWriter := io.MultiWriter(chunkFile)

		_, err = io.CopyN(multiWriter, inputFile, chunkSize)
		if err != nil && err != io.EOF {
			return nil, err
		}

		chunkFile.Close()

		// Calculate hash value
		// hashValue := fmt.Sprintf("%x", hasher.Sum(nil))
		// hashedFileName := hashValue

		// Rename the chunk file with its hash value
		// err = os.Rename(chunkFile.Name())
		// if err != nil {
		// 	return nil, nil, err
		// }

		chunkNames = append(chunkNames, chunkFile.Name())
		// hashValues = append(hashValues, hashValue)

		// chunkNames[0] = "ggggggg"
		// hasher.Reset()
	}

	return chunkNames, nil
}
