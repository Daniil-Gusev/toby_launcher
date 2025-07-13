package file_utils

import (
	"encoding/json"
	"io"
	"os"
	"toby_launcher/apperrors"
)

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, apperrors.New(apperrors.Err, "Error opening file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, apperrors.New(apperrors.Err, "Error read file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	err = file.Close()
	if err != nil {
		return nil, apperrors.New(apperrors.Err, "Error close file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	return data, nil
}

func WriteFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return apperrors.New(apperrors.Err, "Error creating file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	_, err = file.Write(data)
	if err != nil {
		return apperrors.New(apperrors.Err, "Error writing to file $file: $error", map[string]any{
			"file":  filePath,
			"error": err,
		})
	}

	return nil
}

func EncodeData(data any) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, apperrors.New(apperrors.Err, "Data encode error: $error", map[string]any{
			"error": err,
		})
	}
	return bytes, nil
}

func DecodeData(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		return apperrors.New(apperrors.Err, "Data decode error: $error", map[string]any{
			"error": err,
		})
	}
	return nil
}

func LoadData(filePath string, target any) error {
	data, err := ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := DecodeData(data, target); err != nil {
		return err
	}
	return nil
}

func SaveData(filePath string, target any) error {
	bytes, err := EncodeData(target)
	if err != nil {
		return err
	}
	if err := WriteFile(filePath, bytes); err != nil {
		return err
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
