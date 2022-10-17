const PouchDB = require('pouchdb')

const dbCards = new PouchDB('./database/card')

const findAllUnusedCards = async () => {
    const result = await dbCards.allDocs({ include_docs: true })

    return result.rows.map(row => row.doc)
}

const findOneCardById = (id) => {
    return dbCards.get(id)
}

const insertUnusedCard = async (cardId) => {
    return dbCards.post(cardId)
}

const deleteUnusedCard = async (cardId) => {
    return dbCards.remove(cardId)
}

module.exports = {
    findAllUnusedCards,
    findOneCardById,
    insertUnusedCard,
    deleteUnusedCard
}
