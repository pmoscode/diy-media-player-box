package api

import (
	"controller/database"
	"controller/database/schema"
	"controller/utils"
	"log"
)

type AudioBookService struct {
	dbClient    *database.Database
	cardService *CardService
}

func (a *AudioBookService) GetAllAudioBooks() (*[]schema.AudioBook, error) {
	allAudioBooks, dbResult := a.dbClient.GetAllAudioBooks()
	log.Println(dbResult)

	return allAudioBooks, nil
}

func (a *AudioBookService) AddAudioBook(audioBook *schema.AudioBook) (*schema.AudioBook, error) {
	dbResult := a.dbClient.InsertAudioBook(audioBook)
	log.Println(dbResult)

	if audioBook.CardId != nil {
		err := a.cardService.RemoveUnusedCard(audioBook.CardId.ID)
		if err != nil {
			return nil, err
		}
	}

	utils.CreateMediaFolder(audioBook.ID)

	return audioBook, nil
}

/*

 */

func NewAudioBookService() *AudioBookService {
	databaseSingleton, err := database.CreateDatabase(false)
	if err != nil {
		return nil
	}

	return &AudioBookService{
		dbClient:    databaseSingleton,
		cardService: NewCardService(),
	}
}
