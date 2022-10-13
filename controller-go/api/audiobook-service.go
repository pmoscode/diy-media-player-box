package api

import (
	"controller/database"
	"controller/database/schema"
	"controller/utils"
	"log"
	"mime/multipart"
	"path/filepath"
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

func (a *AudioBookService) UpdateAudioBook(id uint, audioBook *schema.AudioBook) error {
	audioBookDb, dbResult := a.dbClient.GetAudioBookById(id)
	log.Println(dbResult)

	if audioBookDb.CardId != nil && audioBook.CardId != nil && audioBookDb.CardId.ID != audioBook.CardId.ID {
		_, err := a.cardService.AddUnusedCard(audioBookDb.CardId.CardId)
		if err != nil {
			return err
		}
	}

	if audioBook.CardId != nil {
		err := a.cardService.RemoveUnusedCard(audioBook.CardId.ID)
		if err != nil {
			return err
		}
	}

	err := utils.MergeAudioBook(audioBookDb, *audioBook)
	if err != nil {
		return err
	}

	dbResult = a.dbClient.UpdateAudioBook(audioBookDb)
	log.Println(dbResult)

	return nil
}

func (a *AudioBookService) DeleteAudioBook(id uint) error {
	audioBookDb, dbResult := a.dbClient.GetAudioBookById(id)
	log.Println(dbResult)

	a.dbClient.DeleteAudioBook(audioBookDb)

	utils.DeleteMediaFolder(id)

	if audioBookDb.CardId != nil {
		_, err := a.cardService.AddUnusedCard(audioBookDb.CardId.CardId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AudioBookService) UploadTracks(id uint, audioFiles []*multipart.FileHeader) error {
	audioBookDb, dbResult := a.dbClient.GetAudioBookById(id)
	log.Println(dbResult)

	for trackNumber, file := range audioFiles {
		mediaPath := utils.GetCompletePathToMediaFolder(audioBookDb.ID)

		err := utils.CopyRequestFileToMediaFolder(mediaPath, file)
		if err != nil {
			continue
		}

		title, length := utils.GetAudioInformation(filepath.Join(mediaPath, file.Filename))
		track := &schema.AudioTrack{
			Track:    uint(trackNumber),
			Title:    title,
			Length:   length,
			FileName: file.Filename,
		}
		audioBookDb.TrackList = append(audioBookDb.TrackList, track)
	}

	dbResult = a.dbClient.UpdateAudioBook(audioBookDb)
	log.Println(dbResult)

	return nil
}

func (a *AudioBookService) DeleteAllTracks(id uint) (*schema.AudioBook, error) {
	audioBookDb, dbResult := a.dbClient.GetAudioBookById(id)
	log.Println(dbResult)

	utils.DeleteMediaFolderContent(audioBookDb)

	audioBookDb.TrackList = make([]*schema.AudioTrack, 0)

	a.dbClient.UpdateAudioBook(audioBookDb)

	return audioBookDb, nil
}

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
