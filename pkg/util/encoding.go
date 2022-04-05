package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func EncodeJsonFile(filePath string, rs interface{}) error {
	var des io.Writer = os.Stdout
	if strings.TrimSpace(filePath) == "" {
		return errors.New("this file is not exist")
	}
	if _, err :=os.Open(filePath); os.IsExist(err) {
		filePath = fmt.Sprintf("%v_%v", filePath, UUID8())
	}
	return json.NewEncoder(des).Encode(rs)
}

func DecodeJsonFile(filePath string, rs interface{}) error {
	var src io.Reader = os.Stdin
	if strings.TrimSpace(filePath) == "" {
		return errors.New("this file is not exist")
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	src = f
	return json.NewDecoder(src).Decode(rs)
}
