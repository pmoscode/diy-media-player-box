const PouchDB = require('pouchdb')

const dbAudioBook = new PouchDB('./database/audio-book')

const findAllAudioBooks = async () => {
    const result = await dbAudioBook.allDocs({ include_docs: true })

    return result.rows.map(row => row.doc)
}

const findOneAudioBook = (audioBook) => {
    return dbAudioBook.get(audioBook.id)
}

const findOneAudioBookById = (id) => {
    return dbAudioBook.get(id)
}

const insertAudioBook = async (audioBook) => {
    return dbAudioBook.post(audioBook)
}

const updateAudioBook = async (audioBook) => {
    return dbAudioBook.put(audioBook)
}

const deleteAudioBook = async (audioBook) => {
    return dbAudioBook.remove(audioBook)
}

module.exports = {
    findAllAudioBooks,
    findOneAudioBook,
    findOneAudioBookById,
    insertAudioBook,
    updateAudioBook,
    deleteAudioBook
}
