/*
An API client for opensubtitles.org

This is a client for the OSDb protocol. Currently the package only allows movie
identification, and subtitles search.
*/
package osdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	ChunkSize = 65536 // 64k
)

func Hash(path string) (hash uint64, err error) {
	// Check file size.
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return 0, fmt.Errorf("File is too small")
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return 0, err
	}
	for _, num := range nums {
		hash += num
	}

	return hash + uint64(fi.Size()), nil
}

// Fill buffer
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("Invalid read", n)
	}
	return
}