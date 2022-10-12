package main

import (
	"controller/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Static("ui", "/")

	apiRoute := router.Group("/api")
	{
		apiRoute.GET("/audio-books", api.GetAllAudioBooks)
		apiRoute.POST("/audio-books", api.AddAudioBook)
		apiRoute.PATCH("/audio-books/:id", api.UpdateAudioBook)
		apiRoute.DELETE("/audio-books/:id", api.DeleteAudioBook)

		apiRoute.POST("/audio-books/:id/tracks", api.UploadTracks)
		apiRoute.DELETE("/audio-books/:id/tracks", api.DeleteAllTracks)

		apiRoute.POST("/audio-books/:uid/play", api.PlayAudioBook)
		apiRoute.POST("/audio-books/:id/track/:track/play", api.PlayTrack)
		apiRoute.POST("/audio-books/pause", api.PauseTrack)
		apiRoute.POST("/audio-books/stop", api.StopTrack)
	}

	router.Run("2020")

	//db, _ := database.CreateDatabase(false)
	//
	//ab := &schema.AudioBook{
	//	Title:       "One",
	//	CardId:      "123456",
	//	TimesPlayed: 2,
	//}
	//
	//track := schema.AudioTrack{
	//	Track:    "1",
	//	Title:    "T1",
	//	Length:   30,
	//	FileName: "/tmp/dart.mp3",
	//}
	//ab.TrackList = append(ab.TrackList, track)
	//
	//db.InsertAudioBook(ab)
	//utils.PrintFormatStruct(ab)
	//
	//ab.Title = "Two"
	//
	//db.UpdateAudioBook(ab)
	//utils.PrintFormatStruct(ab)
	//
	//ab2, _ := db.GetAllAudioBooks()
	//utils.PrintPrettyFormatStruct(ab2)
	//
	//db.DeleteAudioBook(ab)
	//
	//ab3, _ := db.GetAllAudioBooks()
	//utils.PrintPrettyFormatStruct(ab3)
}
