package fileutils

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
)

func ValidateFileName(name string) error {
	if strings.Contains(name, "..") {
		return errors.New("invalid file name")
	}
	return nil
}

func OpenFile(fileName string, folderName string) (string, os.FileInfo, *os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", nil, nil, err
	}
	filePath := path.Join(dir, folderName, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return "", nil, nil, err
	}
	return filePath, stat, file, nil
}

func HandleFileOpenError(w http.ResponseWriter, err error) {
	if os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
