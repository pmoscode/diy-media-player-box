const buildEmptyDbCard = (cardId = '') => {
    return {
        id: '',
        cardId: cardId
    }
}

const mergeCard = (cardSource, cardTarget) => {
    return Object.assign(cardTarget, cardSource)
}

const getCardFromRequest = (req) => {
    if (req && req.body) {
        return mergeCard(req.body, buildEmptyDbCard())
    }
}

const getCardForResponse = (card) => {
    if (card) {
        if (card._id) {
            card.id = card._id
        }
        delete card.ok
        delete card._id
        delete card._rev
        delete card.rev

        return card
    }
}

module.exports = {
    buildEmptyDbCard,
    mergeCard,
    getCardFromRequest,
    getCardForResponse
}
