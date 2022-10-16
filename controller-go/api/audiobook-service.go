package api

import (
	uiSchema "controller/api/schema"
	"controller/database"
	dbSchema "controller/database/schema"
	"controller/utils"
	"mime/multipart"
	"path/filepath"
	"sort"
	"time"
)

type AudioBookService struct {
	dbClient    *database.Database
	cardService *CardService
}

func (a *AudioBookService) GetAllAudioBooks() ([]*uiSchema.AudioBook, error) {
	allAudioBooks, _ := a.dbClient.GetAllAudioBooks()

	audioBooks := make([]*uiSchema.AudioBook, 0)
	for _, audioBook := range *allAudioBooks {
		converted := utils.ConvertAudioBookDbToUi(&audioBook)
		audioBooks = append(audioBooks, converted)
	}

	return audioBooks, nil
}

func (a *AudioBookService) AddAudioBook(audioBook *uiSchema.AudioBook) (*uiSchema.AudioBook, error) {
	dbAudioBook := utils.ConvertAudioBookUiToDb(audioBook)
	dbAudioBook.TimesPlayed = 0
	dbAudioBook.LastPlayed = time.Now()
	a.dbClient.InsertAudioBook(dbAudioBook)

	if audioBook.Card != nil {
		err := a.cardService.RemoveUnusedCard(uint(audioBook.Card.Id))
		if err != nil {
			return nil, err
		}
	}

	utils.CreateMediaFolder(dbAudioBook.ID)

	return audioBook, nil
}

func (a *AudioBookService) UpdateAudioBook(id uint, audioBookUi *uiSchema.AudioBook) error {
	audioBookDb, _ := a.dbClient.GetAudioBookById(id)

	if audioBookDb.CardId != "" && audioBookUi.Card != nil && audioBookDb.CardId != audioBookUi.Card.CardId {
		_, err := a.cardService.AddUnusedCard(audioBookDb.CardId)
		if err != nil {
			return err
		}
	}

	if audioBookUi.Card != nil {
		err := a.cardService.RemoveUnusedCard(uint(audioBookUi.Card.Id))
		if err != nil {
			return err
		}
	}

	utils.MergeAudioBookUiToDb(audioBookDb, audioBookUi)

	a.dbClient.UpdateAudioBook(audioBookDb)

	return nil
}

func (a *AudioBookService) DeleteAudioBook(id uint) (*uiSchema.AudioBook, error) {
	audioBookDb, _ := a.dbClient.GetAudioBookById(id)

	a.dbClient.DeleteAudioBook(audioBookDb)

	utils.DeleteMediaFolder(id)

	if audioBookDb.CardId != "" {
		_, err := a.cardService.AddUnusedCard(audioBookDb.CardId)
		if err != nil {
			return nil, err
		}
	}

	return utils.ConvertAudioBookDbToUi(audioBookDb), nil
}

func (a *AudioBookService) UploadTracks(id uint, audioFiles []*multipart.FileHeader) ([]*uiSchema.AudioTrack, error) {
	audioBookDb, _ := a.dbClient.GetAudioBookById(id)

	uiTracks := make([]*dbSchema.AudioTrack, 0)
	trackSum := len(audioBookDb.TrackList) + 1

	sort.SliceStable(audioFiles, func(i, j int) bool {
		return audioFiles[i].Filename < audioFiles[j].Filename
	})

	for trackNumber, file := range audioFiles {
		mediaPath := utils.GetCompletePathToMediaFolder(audioBookDb.ID)

		err := utils.CopyRequestFileToMediaFolder(mediaPath, file)
		if err != nil {
			continue
		}

		title, length := utils.GetAudioInformation(filepath.Join(mediaPath, file.Filename))
		track := &dbSchema.AudioTrack{
			Track:    uint(trackNumber + trackSum),
			Title:    title,
			Length:   length,
			FileName: file.Filename,
		}
		audioBookDb.TrackList = append(audioBookDb.TrackList, track)
		uiTracks = append(uiTracks, track)
	}

	a.dbClient.UpdateAudioBook(audioBookDb)

	return utils.ConvertAudioBookTracksDbToUi(uiTracks), nil
}

func (a *AudioBookService) DeleteAllTracks(id uint) (*uiSchema.AudioBook, error) {
	audioBookDb, _ := a.dbClient.GetAudioBookById(id)

	utils.DeleteMediaFolderContent(audioBookDb)

	for _, audioTrack := range audioBookDb.TrackList {
		a.dbClient.DeleteAudioTrack(audioTrack)
	}

	audioBookDb.TrackList = make([]*dbSchema.AudioTrack, 0)

	a.dbClient.UpdateAudioBook(audioBookDb)

	return utils.ConvertAudioBookDbToUi(audioBookDb), nil
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
