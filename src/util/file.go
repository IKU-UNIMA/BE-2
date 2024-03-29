package util

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
)

func CheckFileType(file *multipart.FileHeader) error {
	fileType := file.Header["Content-Type"][0]
	switch fileType {
	case "application/pdf":
		return nil
	case "image/png":
		return nil
	case "image/jpeg":
		return nil
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return nil
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return nil
	case "text/plain":
		return nil
	}

	return errors.New(fmt.Sprint("unsupported file type for ", file.Filename))
}

func CheckFileIsPDF(file *multipart.FileHeader) error {
	fileType := file.Header["Content-Type"][0]
	if fileType == "application/pdf" {
		return nil
	}

	return errors.New(fmt.Sprint("unsupported file type for ", file.Filename))
}

func CreateFileUrl(fileId string) string {
	return "https://drive.google.com/file/d/" + fileId + "/view"
}

func GetFileIdFromUrl(url string) string {
	arr := strings.Split(url, "/")
	return arr[len(arr)-2]
}
