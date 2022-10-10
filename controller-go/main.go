package main

import (
	"controller/database"
	"log"
)

func main() {
	//router := gin.Default()
	//
	//router.GET("/audio-books", api.GetAllAudioBooks)
	//router.POST("/audio-books", api.AddAudioBook)
	//router.PATCH("/audio-books/:id", api.UpdateAudioBook)
	//router.DELETE("/audio-books/:id", api.DeleteAudioBook)
	//
	//router.POST("/audio-books/:id/tracks", api.UploadTracks)
	//router.DELETE("/audio-books/:id/tracks", api.DeleteAllTracks)
	//
	//router.POST("/audio-books/:uid/play", api.PlayAudioBook)
	//router.POST("/audio-books/:id/track/:track/play", api.PlayTrack)
	//router.POST("/audio-books/pause", api.PauseTrack)
	//router.POST("/audio-books/stop", api.StopTrack)
	//
	//router.Run("2020")

	db, _ := database.CreateDatabase(false)

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
	//log.Println(ab)
	//
	//ab.Title = "Two"
	//
	//db.UpdateAudioBook(ab)
	//log.Println(ab)

	//db.DeleteAudioBook(ab)

	ab2, res := db.GetAllAudioBooks()
	log.Println(ab2, res)
}
