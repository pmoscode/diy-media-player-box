package api

import (
	uiSchema "controller/api/schema"
	"controller/database"
	dbSchema "controller/database/schema"
	"controller/mqtt"
	"controller/utils"
	"fmt"
	"github.com/pmoscode/go-common/heartbeat"
	mqtt2 "github.com/pmoscode/go-common/mqtt"
	"github.com/pmoscode/go-common/shutdown"
	"github.com/pmoscode/go-common/yamlconfig"
	"log"
	"mime/multipart"
	"path/filepath"
	"sort"
	"time"
)

var mqttClient *mqtt2.Client

type AudioBookService struct {
	dbClient      *database.Database
	cardService   *CardService
	lastPlayedUid string
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

	if (audioBookDb.CardId != "" && audioBookUi.Card == nil) || (audioBookDb.CardId != "" && audioBookUi.Card != nil && audioBookDb.CardId != audioBookUi.Card.CardId) {
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

	message := &mqtt2.Message{
		Topic: "/controller/play",
		Value: request,
	}

	mqttClient.Publish(message)

	return nil
}

func (a *AudioBookService) StopAudioTrack() error {
	message := &mqtt2.Message{
		Topic: "/controller/stop",
		Value: nil,
	}

	mqttClient.Publish(message)

	return nil
}

func (a *AudioBookService) PauseAudioTrack() error {
	// TODO Switch is removed. Has to be solved somehow.
	message := &mqtt2.Message{
		Topic: "/controller/switch",
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

	heartBeat := heartbeat.New(10*time.Second, sendHeartbeat)
	heartBeat.Run()

	var config Config
	err = yamlconfig.LoadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("Could not load config file")
	}

	service := &AudioBookService{
		dbClient:      databaseSingleton,
		cardService:   NewCardService(),
		lastPlayedUid: "",
	}

	mqttClient = mqtt2.NewClient(mqtt2.WithBroker(config.MqttBroker.Host, 1883),
		mqtt2.WithClientId(config.Controller.MqttClientId),
		mqtt2.WithOrderMatters(false))
	err = mqttClient.Connect()
	if err != nil {
		if err != nil {
			log.Fatalln("MQTT broker not found... exiting.")
		}
	}
	shutdown.GetObserver().AddCommand(func() error {
		err2 := mqttClient.Disconnect()
		if err2 != nil {
			return err2
		}

		return nil
	})

	mqttClient.Subscribe("/rfid-reader/cardId", service.OnMessageReceivedCardId)
	mqttClient.Subscribe("/audio-player/done", service.OnMessageReceivedPlayDone)

	return service
}

func (a *AudioBookService) OnMessageReceivedCardId(message mqtt2.Message) {
	card := &mqtt.RfidReaderSubscribeMessage{}
	message.ToStruct(card)

	sendStatusMessage(mqtt2.Info, "Got card id: ", card.CardId)

	audioPlayerMessage := &mqtt2.Message{}

	if card.CardId == "" {
		if a.lastPlayedUid == "" {
			sendStatusMessage(mqtt2.Info, "Card removed, but nothing is played currently...")
		} else {
			audioPlayerMessage.Topic = "/controller/pause"
			audioPlayerMessage.Value = nil
			mqttClient.Publish(audioPlayerMessage)

			sendStatusMessage(mqtt2.Info, "Going for a short pause...")
		}
	} else {
		if a.lastPlayedUid != card.CardId {
			audioBookDb, dbResult := a.dbClient.GetAudioBookByCardId(card.CardId)

			if dbResult != database.DbRecordNotFound {
				sendStatusMessage(mqtt2.Info, "Going to play new audio book...")

				request := &mqtt.AudioPlayerPublishMessage{
					Id:        audioBookDb.ID,
					TrackList: []string{},
				}

				for _, track := range audioBookDb.TrackList {
					mediaPath := utils.GetCompletePathToMediaFolder(audioBookDb.ID)
					audioFilePath := filepath.Join(mediaPath, track.FileName)

					request.TrackList = append(request.TrackList, audioFilePath)
				}
				a.lastPlayedUid = card.CardId
				sendStatusMessage(mqtt2.Info, "Setting lastPlayedUid to: ", a.lastPlayedUid)

				audioPlayerMessage.Topic = "/controller/play"
				audioPlayerMessage.Value = request
				mqttClient.Publish(audioPlayerMessage)

				sendStatusMessage(mqtt2.Info, "Playing now: ", audioBookDb.Title)
			} else {
				_, dbResult := a.dbClient.GetCard(card.CardId)

				if dbResult == database.DbRecordNotFound {
					a.dbClient.AddUnusedCard(card.CardId)

					sendStatusMessage(mqtt2.Info, "Added new card: ", card.CardId)
				} else {
					sendStatusMessage(mqtt2.Warn, "Card not assigned: ", card.CardId)
				}
			}
		} else {
			audioPlayerMessage.Topic = "/controller/resume"
			audioPlayerMessage.Value = nil

			mqttClient.Publish(audioPlayerMessage)
		}

	}
}

func (a *AudioBookService) OnMessageReceivedPlayDone(message mqtt2.Message) {
	playDone := &mqtt.PlayDoneSubscribeMessage{}
	message.ToStruct(playDone)

	audioBookDb, dbResult := a.dbClient.GetAudioBookById(playDone.Uid)

	if dbResult != database.DbRecordNotFound {
		audioBookDb.LastPlayed = time.Now()
		audioBookDb.TimesPlayed++

		a.dbClient.UpdateAudioBook(audioBookDb)
	}

	sendStatusMessage(mqtt2.Info, "Play done -> ", playDone.Uid)
	a.lastPlayedUid = ""
}

func sendHeartbeat() {
	mqttClient.Publish(&mqtt2.Message{
		Topic: "/heartbeat/controller",
		Value: &mqtt2.StatusPublishMessage{
			Status: "online",
			Type:   mqtt2.Info,
		},
	})
}

func sendStatusMessage(messageType mqtt2.StatusType, message ...any) {
	messageTxt := fmt.Sprint(message...)

	mqttMessage := &mqtt2.StatusPublishMessage{
		Type:      messageType,
		Status:    messageTxt,
		Timestamp: time.Now(),
	}

	mqttClient.Publish(&mqtt2.Message{
		Topic: "/status/controller",
		Value: mqttMessage,
	})
}
