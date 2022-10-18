package api

import (
	"controller/api/schema"
	"controller/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

var audioBookService = NewAudioBookService()
var cardService = NewCardService()

// GetAllAudioBooks godoc
// @Summary         GetAllAudioBooks
// @Description     Get all audiobooks stored
// @Tags            audio-book
// @Produce         json
// @Success         200  {object}  schema.AudioBookFull
// @Failure         500  {string}  string "Internal Server Error"
// @Router          /audio-books [get]
func GetAllAudioBooks(c echo.Context) error {
	audioBooks, err := audioBookService.GetAllAudioBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, audioBooks)
}

// AddAudioBook  godoc
// @Summary      AddAudioBook
// @Description  Add a new audiobook
// @Tags         audio-book
// @Accept		 json
// @Produce      json
// @Success      200  {object}  schema.AudioBookFull
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /audio-books [post]
func AddAudioBook(c echo.Context) error {
	var audioBook schema.AudioBookUi

	err := c.Bind(&audioBook)
	if err == nil {
		audioBookResult, err := audioBookService.AddAudioBook(&audioBook)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
		}
		return c.JSON(http.StatusOK, audioBookResult)
	}

	return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
}

// UpdateAudioBook godoc
// @Summary        UpdateAudioBook
// @Description    Update an existing audiobook
// @Tags           audio-book
// @Accept		   json
// @Produce        json
// @Param          id           path   uint                 true   "id of audio-book"
// @Param          audio-book   body   schema.AudioBookUi   true   "content of audio-book"
// @Success        200  {object}  schema.AudioBookFull
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /audio-books/{id} [patch]
func UpdateAudioBook(c echo.Context) error {
	var audioBookUi schema.AudioBookUi

	id := c.Param("id")

	errBind := c.Bind(&audioBookUi)
	if errBind == nil {
		audioBook, err := audioBookService.UpdateAudioBook(utils.ConvertToUint(id), &audioBookUi)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
		}

		return c.JSON(http.StatusOK, &audioBook)
	}

	return c.JSON(http.StatusInternalServerError, &response{message: errBind.Error()})
}

// DeleteAudioBook godoc
// @Summary        DeleteAudioBook
// @Description    Delete an existing audiobook
// @Tags           audio-book
// @Produce        json
// @Param          id   path   uint   true   "id of audio-book"
// @Success        200  {object}  schema.AudioBookFull
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /audio-books/{id} [delete]
func DeleteAudioBook(c echo.Context) error {
	id := c.Param("id")

	audioBook, err := audioBookService.DeleteAudioBook(utils.ConvertToUint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, audioBook)
}

// UploadTracks godoc
// @Summary        UploadTracks
// @Description    Uploads audio tracks to an existing audiobook
// @Tags           audio-track
// @Accept		   json
// @Produce        json
// @Param          id       path       uint                true   "id of audio-book"
// @Param          tracks   formData   schema.AudioTrack   true   "Audio files to upload"
// @Success        200  {object}  schema.AudioBookFull
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /audio-books/{id}/tracks [post]
func UploadTracks(c echo.Context) error {
	id := c.Param("id")
	files := utils.GetAllFilesFromRequest(c)

	tracks, err := audioBookService.UploadTracks(utils.ConvertToUint(id), files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, tracks)
}

// DeleteAllTracks godoc
// @Summary        DeleteAllTracks
// @Description    Delete all audio tracks of an existing audiobook
// @Tags           audio-track
// @Produce        json
// @Param          id   path   uint   true   "id of audio-book"
// @Success        200  {object}  schema.AudioBookFull
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /audio-books/{id}/tracks [delete]
func DeleteAllTracks(c echo.Context) error {
	id := c.Param("id")

	audioBook, err := audioBookService.DeleteAllTracks(utils.ConvertToUint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, audioBook)
}

// GetAllUnassignedCards godoc
// @Summary        GetAllUnassignedCards
// @Description    Get all unassigned cards (rfid card ids)
// @Tags           card
// @Produce        json
// @Success        200  {object}  schema.Card
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /cards/unassigned [get]
func GetAllUnassignedCards(c echo.Context) error {
	cards, err := cardService.GetAllUnusedCards()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.JSON(http.StatusOK, cards)
}

// PlayTrack godoc
// @Summary        PlayTrack
// @Description    Plays an audio tracks of an existing audiobook
// @Tags           audio-book-debug
// @Produce        json
// @Param          id      path   uint   true   "id of audio-book"
// @Param          track   path   uint   true   "track number to be played"
// @Success        200  {string}  string "No Content"
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /audio-books/{id}/track/{track}/play [post]
func PlayTrack(c echo.Context) error {
	id := c.Param("id")
	idxTrack := c.Param("track")

	err := audioBookService.PlayAudioTrack(utils.ConvertToUint(id), utils.ConvertToUint(idxTrack))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// PauseTrack godoc
// @Summary        PauseTrack
// @Description    Pauses the current playing audio track (if any) - Can be called again to resume playback
// @Tags           audio-book-debug
// @Produce        json
// @Success        200  {string}  string "No Content"
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /audio-books/pause [post]
func PauseTrack(c echo.Context) error {
	err := audioBookService.PauseAudioTrack()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// StopTrack godoc
// @Summary        StopTrack
// @Description    Stops the current playing audio track (if any)
// @Tags           audio-book-debug
// @Produce        json
// @Success        200  {string}  string "No Content"
// @Failure        500  {string}  string "Internal Server Error"
// @Router         /audio-books/stop [post]
func StopTrack(c echo.Context) error {
	err := audioBookService.StopAudioTrack()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response{message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
