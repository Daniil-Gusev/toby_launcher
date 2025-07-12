package main

import (
	"bytes"
	"fmt"
	"github.com/bodgit/sevenzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractFile(f *sevenzip.File, destDir string, perms os.FileMode) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", f.FileHeader.Name, err)
	}
	defer rc.Close()
	destPath := filepath.Join(destDir, f.FileHeader.Name)
	if f.FileHeader.Mode().IsDir() {
		return os.MkdirAll(destPath, perms)
	}
	if strings.HasSuffix(f.FileHeader.Name, ".7z") {
		archive, err := io.ReadAll(rc)
		if err != nil {
			return fmt.Errorf("Failed to read embedded archive: %v", err)
		}
		return extractArchive(archive, destDir, perms)
	}
	if err = writeData(rc, destPath, perms); err != nil {
		return err
	}
	fmt.Printf("Extracted file: %s\n", destPath)
	return nil
}

func extractArchive(data []byte, destDir string, perms os.FileMode) error {
	r, err := sevenzip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to read archive: %v", err)
	}
	for _, f := range r.File {
		if err := extractFile(f, destDir, perms); err != nil {
			return err
		}
	}
	return nil
}
