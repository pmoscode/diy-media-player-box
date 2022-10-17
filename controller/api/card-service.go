package api

import (
	uiSchema "controller/api/schema"
	"controller/database"
	"controller/utils"
)

type CardService struct {
	dbClient *database.Database
}

func (c *CardService) GetAllUnusedCards() ([]*uiSchema.Card, error) {
	allCards, _ := c.dbClient.GetAllCards()

	cards := make([]*uiSchema.Card, 0)
	for _, card := range *allCards {
		cards = append(cards, utils.ConvertCardDbToUi(&card))
	}

	return cards, nil
}

func (c *CardService) AddUnusedCard(cardId string) (*uiSchema.Card, error) {
	card, _ := c.dbClient.AddUnusedCard(cardId)

	return utils.ConvertCardDbToUi(card), nil
}

func (c *CardService) RemoveUnusedCard(id uint) error {
	c.dbClient.RemoveUnusedCard(id)

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
