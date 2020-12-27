const audioBookDbClient = require('./audio-book-db-client')
const cardDbClient = require('./card-db-client')
const audioPlayer = require('../audio-player/audio-player-controller')
const ad = require('./audioDefinition')
const path = require('path')
const logger = require('../logger')(module)

const getFilePathForUid = async (uid) => {
    const audioBook = await findAudioBookWithUid(uid)

    if (!audioBook) {
        return insertUnusedCard(uid)
    }

    if (audioBook.trackList) {
        const firstTrack = audioBook.trackList[0]
        const completePath = getCompletePathToAudioFile(audioBook, firstTrack)

        logger.info('Found ' + audioBook.trackList.length + ' tracks to play...')

        return { filePath: completePath, uid: uid, audioBook: audioBook }
    }

    playPredefinedAudio(ad.audioDefinition.NO_TRACKS_DEFINED)

    logger.info('Nothing found')

    return { filePath: null, uid: uid, audioBook: audioBook }
}

const insertUnusedCard = async (uid) => {
    if (!await checkIfUnusedCardAlreadyInserted(uid)) {
        const card = {
            id: '',
            cardId: uid
        }
        await cardDbClient.insertUnusedCard(card)
        logger.info('Inserted unused card: ' + JSON.stringify(card))
        playPredefinedAudio(ad.audioDefinition.CARD_NOT_ALLOCATED)

        return { filePath: null, uid: uid, audioBook: null }
    }

    logger.info('Unused card, but already stored')

    return { filePath: null, uid: uid }
}

const getProgress = (() => {
    return () => {
        return audioPlayer.getProgress()
    }
})()

const playPredefinedAudio = (audioDefinition, playDelay) => {
    logger.debug(audioDefinition)
}

const checkIfUnusedCardAlreadyInserted = async (uid) => {
    const allCards = await cardDbClient.findAllUnusedCards()
    const unusedCard = allCards.filter(card => card.cardId === uid)

    return unusedCard.length !== 0
}

const findAudioBookWithUid = async (uid) => {
    const allAudioBooks = await audioBookDbClient.findAllAudioBooks()
    const audioBook = allAudioBooks
        .filter(audioBook => audioBook.card)
        .filter(audioBook => audioBook.card.cardId === uid)

    return (audioBook && audioBook.length > 0) ? audioBook[0] : undefined
}

const getCompletePathToAudioFile = (audioBook, trackList) => {
    return path.join(process.cwd(), 'media', audioBook._id, trackList.fileName)
}

const updateMetadata = (audioBook) => {
    audioBook.timesPlayed++
    audioBook.lastPlayed = new Date()

    return audioBookDbClient.updateAudioBook(audioBook)
}

module.exports = {
    getProgress,
    getFilePathForUid,
    getCompletePathToAudioFile,
    updateMetadata
}
