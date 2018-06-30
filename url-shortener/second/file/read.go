package file

import (
	"os"
	"bytes"
	"io"
)

func Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
