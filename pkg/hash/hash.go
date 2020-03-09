package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func Verify(path string, hash []byte) error {
	hasher := sha256.New()
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}
	if !bytes.Equal(hash, hasher.Sum(nil)) {
		return fmt.Errorf("hash does not match file")
	}
	return nil
}
