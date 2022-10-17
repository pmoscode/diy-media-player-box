package utils

import (
	"github.com/labstack/echo/v4"
	"mime/multipart"
)

func GetAllFilesFromRequest(c echo.Context) []*multipart.FileHeader {
	var result = make([]*multipart.FileHeader, 0)

	multipartForm, err2 := c.MultipartForm()
	if err2 != nil {
		return nil
	}

	for key, _ := range multipartForm.File {
		file, err := c.FormFile(key)
		if err != nil {
			return nil
		}
		result = append(result, file)
	}

	return result
}
