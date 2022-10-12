package api

import (
	"controller/database"
	"controller/database/schema"
	"log"
)

type CardService struct {
	dbClient *database.Database
}

func (c *CardService) GetAllUnusedCards() (*[]schema.Card, error) {
	allCards, dbResult := c.dbClient.GetAllCards()
	log.Println(dbResult)

	return allCards, nil
}

func (c *CardService) AddUnusedCard(cardId string) (*schema.Card, error) {
	card, dbResult := c.dbClient.AddUnusedCard(cardId)
	log.Println(dbResult)

	return card, nil
}

func (c *CardService) RemoveUnusedCard(id uint) error {
	dbResult := c.dbClient.RemoveUnusedCard(id)
	log.Println(dbResult)

	return nil
}

func NewCardService() *CardService {
	databaseSingleton, err := database.CreateDatabase(false)
	if err != nil {
		return nil
	}

	return &CardService{
		dbClient: databaseSingleton,
	}
}
