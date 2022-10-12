package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

func CreateMediaFolder(id uint) bool {
	mediaFolder := GetCompletePathToMediaFolder(id)
	err := os.MkdirAll(mediaFolder, os.ModePerm)
	if err != nil {
		return false
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
