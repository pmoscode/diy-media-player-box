package utils

import (
	"controller/database/schema"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

func CreateMediaFolder(id uint) bool {
	mediaFolder := GetCompletePathToMediaFolder(id)
	err := os.MkdirAll(mediaFolder, os.ModePerm)

	return err == nil
}

func DeleteMediaFolder(id uint) bool {
	mediaFolder := GetCompletePathToMediaFolder(id)
	err := os.RemoveAll(mediaFolder)

	return err == nil
}

func DeleteMediaFolderContent(audioBook *schema.AudioBook) bool {
	for _, file := range audioBook.TrackList {
		filePath := filepath.Join(GetCompletePathToMediaFolder(audioBook.ID), file.FileName)
		err := os.Remove(filePath)
		if err != nil {
			return false
		}
	}
	return true
}

func GetCompletePathToMediaFolder(id uint) string {
	return filepath.Join(GetPathOfCurrentBinary(), "media", strconv.Itoa(int(id)))
}

func GetPathOfCurrentBinary() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func CopyRequestFileToMediaFolder(targetFolder string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(filepath.Join(targetFolder, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
