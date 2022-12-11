package database

import (
	"controller/database/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"os"
	time2 "time"
)

type DbResult int

const (
	DbOk DbResult = iota
	DbRecordNotFound
	DbNoRowsAffected
	DbError
)

var databaseSingleton *Database = nil

type Database struct {
	db *gorm.DB
}

func (r *Database) initDatabase(dbFilename string, debug bool) {
	logLevel := logger.Silent
	if debug {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time2.Second, // Slow SQL threshold
			LogLevel:                  logLevel,     // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(dbFilename+".db"), &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: true,
	})
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&schema.AudioBook{}, schema.AudioTrack{}, schema.Card{})
	if err != nil {
		log.Println(err)
		panic("Could not migrate schema for 'AudioBook'")
	}

	r.db = db
}

func (r *Database) GetAllAudioBooks() (*[]schema.AudioBook, DbResult) {
	var data []schema.AudioBook

	result := r.db.Preload(clause.Associations).Find(&data)

	if result.RowsAffected == 0 {
		return &data, DbRecordNotFound
	}

	return &data, DbOk
}

func (r *Database) GetAudioBook(audioBook *schema.AudioBook) (*schema.AudioBook, DbResult) {
	return r.GetAudioBookById(audioBook.ID)
}

func (r *Database) GetAudioBookById(id uint) (*schema.AudioBook, DbResult) {
	var data schema.AudioBook

	result := r.db.Preload(clause.Associations).Where(&schema.AudioBook{
		Model: gorm.Model{
			ID: id,
		},
	}).First(&data)

	if result.RowsAffected == 0 {
		return &data, DbRecordNotFound
	}

	return &data, DbOk
}

func (r *Database) GetAudioBookByCardId(cardId string) (*schema.AudioBook, DbResult) {
	var data schema.AudioBook

	result := r.db.Preload(clause.Associations).Where(&schema.AudioBook{
		CardId: cardId,
	}).First(&data)

	if result.RowsAffected == 0 {
		return &data, DbRecordNotFound
	}

	return &data, DbOk
}

func (r *Database) InsertAudioBook(audioBook *schema.AudioBook) DbResult {
	result := r.db.Create(&audioBook)

	if result.RowsAffected == 0 {
		return DbNoRowsAffected
	}

	return DbOk
}

func (r *Database) UpdateAudioBook(audioBook *schema.AudioBook) DbResult {
	result := r.db.Preload(clause.Associations).Save(&audioBook)

	if result.RowsAffected == 0 {
		return DbNoRowsAffected
	}

	return DbOk
}

func (r *Database) DeleteAudioBook(audioBook *schema.AudioBook) DbResult {
	result := r.db.Unscoped().Select(clause.Associations).Delete(&audioBook)

	if result.RowsAffected == 0 {
		return DbNoRowsAffected
	}

	return DbOk
}

func (r *Database) DeleteAudioTrack(audioTrack *schema.AudioTrack) DbResult {
	result := r.db.Unscoped().Delete(&audioTrack)

	if result.RowsAffected == 0 {
		return DbNoRowsAffected
	}

	return DbOk
}

func (r *Database) GetAllCards() (*[]schema.Card, DbResult) {
	var data []schema.Card

	result := r.db.Find(&data)

	if result.RowsAffected == 0 {
		return &data, DbRecordNotFound
	}

	return &data, DbOk
}

func (r *Database) GetCard(cardId string) (*schema.Card, DbResult) {
	var data schema.Card

	query := schema.Card{
		CardId: cardId,
	}

	result := r.db.Where(query).First(&data)

	if result.RowsAffected == 0 {
		return &data, DbRecordNotFound
	}

	return &data, DbOk
}

func (r *Database) AddUnusedCard(cardId string) (*schema.Card, DbResult) {
	data := schema.Card{
		CardId: cardId,
	}

	result := r.db.Save(&data)

	if result.RowsAffected == 0 {
		return &data, DbNoRowsAffected
	}

	return &data, DbOk
}

func (r *Database) RemoveUnusedCard(id uint) DbResult {
	data := schema.Card{
		Model: gorm.Model{
			ID: id,
		},
	}

	result := r.db.Unscoped().Delete(&data)

	if result.RowsAffected == 0 {
		return DbNoRowsAffected
	}

	return DbOk
}

func CreateDatabase(debug bool) (*Database, error) {
	if databaseSingleton == nil {
		databaseSingleton = &Database{}
		databaseSingleton.initDatabase("data", debug)
	}

	return databaseSingleton, nil
}
