package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func InitDir() error {
	err := os.MkdirAll("archives", 0755)
	if err != nil {
		return fmt.Errorf("error creating archives directory: %w", err)
	}
	return nil
}

func CreateZip(title string) (*zip.Writer, *os.File, error) {
	fullPath := filepath.Join("archives", title)
	file, err := os.Create(fullPath + ".zip")
	if err != nil {
		return nil, nil, fmt.Errorf("error creating archives zip file: %w", err)
	}
	zipWriter := zip.NewWriter(file)
	return zipWriter, file, nil
}

func WriteZip(Writer *zip.Writer, body []byte, url string) error {
	header := &zip.FileHeader{
		Name: url,
	}
	file, err := Writer.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("error creating header zip file: %w", err)
	}
	_, err = io.Copy(file, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error copying body to zip file: %w", err)
	}
	return nil
}
