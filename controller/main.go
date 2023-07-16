package main

import (
	"controller/api"
	_ "controller/docs"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

// @title DIY Music Box for Children
// @version 1.0
// @description This is the controller app for the DIY Music Box for Children.

// @contact.name pmoscode
// @contact.url https://pmoscode.de
// @contact.email info@pmoscode.de

// @license.name GNU General Public License v3.0
// @license.url https://github.com/pmoscode/diy-media-player-box/-/raw/master/LICENSE

// @host localhost:2020
// @BasePath /api
func main() {
	router := echo.New()
	router.Use(middleware.Recover())

	router.Static("/", "ui")
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	apiRoute := router.Group("/api")
	{
		apiRoute.GET("/audio-books", api.GetAllAudioBooks)
		apiRoute.POST("/audio-books", api.AddAudioBook)
		apiRoute.PATCH("/audio-books/:id", api.UpdateAudioBook)
		apiRoute.DELETE("/audio-books/:id", api.DeleteAudioBook)

		apiRoute.POST("/audio-books/:id/tracks", api.UploadTracks)
		apiRoute.DELETE("/audio-books/:id/tracks", api.DeleteAllTracks)

		apiRoute.POST("/audio-books/:id/track/:track/play", api.PlayTrack)
		apiRoute.POST("/audio-books/pause", api.PauseTrack)
		apiRoute.POST("/audio-books/stop", api.StopTrack)

		apiRoute.GET("/cards/unassigned", api.GetAllUnassignedCards)
	}

	for _, route := range router.Routes() {
		log.Println(fmt.Sprintf("%-6s", route.Method), " - ", route.Path)
	}

	router.Logger.Fatal(router.Start(":2020"))
}
