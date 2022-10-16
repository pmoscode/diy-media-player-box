package api

import (
	"controller/api/schema"
	"controller/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

var audioBookService = NewAudioBookService()
var cardService = NewCardService()

func GetAllAudioBooks(c echo.Context) error {
	audioBooks, err := audioBookService.GetAllAudioBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, audioBooks)
}

func AddAudioBook(c echo.Context) error {
	var audioBook schema.AudioBook

	err := c.Bind(&audioBook)
	if err == nil {
		audioBookResult, err := audioBookService.AddAudioBook(&audioBook)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
		}
		return c.JSON(http.StatusOK, audioBookResult)
	}

	return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
}

func UpdateAudioBook(c echo.Context) error {
	var audioBook schema.AudioBook

	id := c.Param("id")

	err := c.Bind(&audioBook)
	if err == nil {
		err = audioBookService.UpdateAudioBook(utils.ConvertToUint(id), &audioBook)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
		}
		return c.JSON(http.StatusOK, &audioBook)
	}

	return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
}

func DeleteAudioBook(c echo.Context) error {
	id := c.Param("id")

	audioBook, err := audioBookService.DeleteAudioBook(utils.ConvertToUint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, audioBook)
}

func UploadTracks(c echo.Context) error {
	id := c.Param("id")
	files := utils.GetAllFilesFromRequest(c)

	tracks, err := audioBookService.UploadTracks(utils.ConvertToUint(id), files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, tracks)
}

func DeleteAllTracks(c echo.Context) error {
	id := c.Param("id")

	audioBook, err := audioBookService.DeleteAllTracks(utils.ConvertToUint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, audioBook)

}

func GetAllUnassignedCards(c echo.Context) error {
	cards, err := cardService.GetAllUnusedCards()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, cards)
}

func PlayAudioBook(c echo.Context) error {

	return nil
}

func PlayTrack(c echo.Context) error {

	return nil
}

func PauseTrack(c echo.Context) error {

	return nil
}

func StopTrack(c echo.Context) error {

	return nil
}
