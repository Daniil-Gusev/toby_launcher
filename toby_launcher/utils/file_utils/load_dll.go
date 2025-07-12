//go:build windows

package file_utils

import (
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"toby_launcher/apperrors"
)

func LoadDLL(libName string) (*windows.DLL, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)
	libDir := filepath.Join(exeDir, "lib")
	var libPath string
	if path := filepath.Join(exeDir, libName); Exists(path) {
		libPath = path
	} else if path := filepath.Join(libDir, libName); Exists(path) {
		libPath = path
	}
	if libPath == "" {
		return nil, apperrors.New(apperrors.Err, "Error loading DLL: library $lib not found.", map[string]any{"lib": libName})
	}
	dll, err := windows.LoadDLL(libPath)
	if err != nil {
		return nil, apperrors.New(apperrors.Err, "Error loading DLL: $error", map[string]any{"error": err.Error()})
	}
	return dll, nil
}
