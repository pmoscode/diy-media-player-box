const cardDbClient = require('../../service/card-db-client')
const cardHelper = require('./card-helper')

const getAllUnassignedCards = async () => {
    const completeResult = await cardDbClient.findAllUnusedCards()

    return completeResult.map(doc => cardHelper.getCardForResponse(doc))
}

const addUnusedCard = async (cardId) => {
    await cardDbClient.insertUnusedCard(cardId)
}

const removeUnusedCard = async (id) => {
    const card = await cardDbClient.findOneCardById(id)

    return cardDbClient.deleteUnusedCard(card)
}

module.exports = {
    getAllUnassignedCards,
    addUnusedCard,
    removeUnusedCard
}
