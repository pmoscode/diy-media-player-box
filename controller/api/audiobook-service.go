package api

import (
	uiSchema "controller/api/schema"
	"controller/database"
	dbSchema "controller/database/schema"
	"controller/mqtt"
	"controller/utils"
	"log"
	"mime/multipart"
	"path/filepath"
	"sort"
	"time"
)

var mqttClient *mqtt.Client

type AudioBookService struct {
	dbClient    *database.Database
	cardService *CardService
}

func (a *AudioBookService) GetAllAudioBooks() ([]*uiSchema.AudioBookFull, error) {
	allAudioBooks, _ := a.dbClient.GetAllAudioBooks()

	audioBooks := make([]*uiSchema.AudioBookFull, 0)
	for _, audioBook := range *allAudioBooks {
		converted := utils.ConvertAudioBookDbToUi(&audioBook)
		audioBooks = append(audioBooks, converted)
	}

	return audioBooks, nil
}

func (a *AudioBookService) AddAudioBook(audioBook *uiSchema.AudioBookUi) (*uiSchema.AudioBookFull, error) {
	audioBookDb := utils.ConvertAudioBookUiToDb(audioBook)
	audioBookDb.TimesPlayed = 0
	audioBookDb.LastPlayed = time.Now()
	a.dbClient.InsertAudioBook(audioBookDb)

	if audioBook.Card != nil {
		err := a.cardService.RemoveUnusedCard(uint(audioBook.Card.Id))
		if err != nil {
			return nil, err
		}
	}

	utils.CreateMediaFolder(audioBookDb.ID)

	return utils.ConvertAudioBookDbToUi(audioBookDb), nil
}

func (a *AudioBookService) UpdateAudioBook(id uint, audioBookUi *uiSchema.AudioBookUi) (*uiSchema.AudioBookFull, error) {
	audioBookDb, _ := a.dbClient.GetAudioBookById(id)

	if audioBookDb.CardId != "" && audioBookUi.Card != nil && audioBookDb.CardId != audioBookUi.Card.CardId {
		_, err := a.cardService.AddUnusedCard(audioBookDb.CardId)
		if err != nil {
			return nil, err
		}
	}

	if audioBookUi.Card != nil {
		err := a.cardService.RemoveUnusedCard(uint(audioBookUi.Card.Id))
		if err != nil {
			return nil, err
		}
	}

	utils.MergeAudioBookUiToDb(audioBookDb, audioBookUi)

	a.dbClient.UpdateAudioBook(audioBookDb)

	return utils.ConvertAudioBookDbToUi(audioBookDb), nil
}

func (a *AudioBookService) DeleteAudioBook(id uint) (*uiSchema.AudioBookFull, error) {
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

func (a *AudioBookService) DeleteAllTracks(id uint) (*uiSchema.AudioBookFull, error) {
	audioBookDb, _ := a.dbClient.GetAudioBookById(id)

	utils.DeleteMediaFolderContent(audioBookDb)

	for _, audioTrack := range audioBookDb.TrackList {
		a.dbClient.DeleteAudioTrack(audioTrack)
	}

	audioBookDb.TrackList = make([]*dbSchema.AudioTrack, 0)

	a.dbClient.UpdateAudioBook(audioBookDb)

	return utils.ConvertAudioBookDbToUi(audioBookDb), nil
}

func (a *AudioBookService) PlayAudioTrack(id uint, idxTrack uint) error {
	audioBookDb, _ := a.dbClient.GetAudioBookById(id)
	track := audioBookDb.TrackList[idxTrack-1]

	mediaPath := utils.GetCompletePathToMediaFolder(id)
	audioFilePath := filepath.Join(mediaPath, track.FileName)

	request := &mqtt.AudioPlayerPublishMessage{
		Id:        id,
		TrackList: []string{audioFilePath},
	}

	message := &mqtt.Message{
		Topic: "/audio-player/play",
		Value: request,
	}

	mqttClient.Publish(message)

	return nil
}

func (a *AudioBookService) StopAudioTrack() error {
	message := &mqtt.Message{
		Topic: "/audio-player/stop",
		Value: nil,
	}

	mqttClient.Publish(message)

	return nil
}

func (a *AudioBookService) PauseAudioTrack() error {
	message := &mqtt.Message{
		Topic: "/audio-player/switch",
		Value: nil,
	}

	mqttClient.Publish(message)

	return nil
}

func NewAudioBookService() *AudioBookService {
	databaseSingleton, err := database.CreateDatabase(false)
	if err != nil {
		return nil
	}

	audioBookService := &AudioBookService{
		dbClient:    databaseSingleton,
		cardService: NewCardService(),
	}

	mqttClient = mqtt.CreateClient(*cliOptions.mqttBrokerIp, 1883, *cliOptions.mqttClientId)
	err = mqttClient.Connect()
	if err != nil {
		if err != nil {
			log.Fatalln("MQTT broker not found... exiting.")
		}
	}
	mqttClient.Subscribe("/rfid-reader/cardId", audioBookService.OnMessageReceivedCardId)
	mqttClient.Subscribe("/audio-player/done", audioBookService.OnMessageReceivedPlayDone)

	return audioBookService
}

func (a *AudioBookService) OnMessageReceivedCardId(message mqtt.Message) {
	card := &mqtt.RfidReaderSubscribeMessage{}
	message.ToStruct(card)

	audioPlayerMessage := &mqtt.Message{}

	if card.CardId == "" {
		audioPlayerMessage.Topic = "/audio-player/switch"
		audioPlayerMessage.Value = nil
	} else {
		audioBookDb, dbResult := a.dbClient.GetAudioBookByCardId(card.CardId)

		if dbResult != database.DbRecordNotFound {
			request := &mqtt.AudioPlayerPublishMessage{
				Id:        audioBookDb.ID,
				TrackList: []string{},
			}

			for _, track := range audioBookDb.TrackList {
				mediaPath := utils.GetCompletePathToMediaFolder(audioBookDb.ID)
				audioFilePath := filepath.Join(mediaPath, track.FileName)

				request.TrackList = append(request.TrackList, audioFilePath)

			}
			audioPlayerMessage.Topic = "/audio-player/play"
			audioPlayerMessage.Value = request
		} else {
			_, dbResult := a.dbClient.GetCard(card.CardId)
			audioPlayerMessage.Topic = "/status/controller"

			if dbResult == database.DbRecordNotFound {
				a.dbClient.AddUnusedCard(card.CardId)

				statusMessage := &mqtt.StatusPublishMessage{
					Type:   mqtt.Info,
					Status: "Added new card: " + card.CardId,
				}

				audioPlayerMessage.Value = statusMessage
			} else {
				statusMessage := &mqtt.StatusPublishMessage{
					Type:   mqtt.Info,
					Status: "Card not assigned: " + card.CardId,
				}

				audioPlayerMessage.Value = statusMessage
			}
		}
	}

	mqttClient.Publish(audioPlayerMessage)
}

func (a *AudioBookService) OnMessageReceivedPlayDone(message mqtt.Message) {
	playDone := &mqtt.PlayDoneSubscribeMessage{}
	message.ToStruct(playDone)

	audioBookDb, dbResult := a.dbClient.GetAudioBookById(playDone.Uid)

	if dbResult != database.DbRecordNotFound {
		audioBookDb.LastPlayed = time.Now()
		audioBookDb.TimesPlayed++

		a.dbClient.UpdateAudioBook(audioBookDb)
	}
}
