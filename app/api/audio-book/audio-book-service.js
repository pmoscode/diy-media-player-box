const path = require('path')
const fs = require('fs')
const logger = require('../../logger')(module)
const audioBookDbClient = require('../../service/audio-book-db-client')
const audioBookHelper = require('./audio-book-helper')
const audioBookPlayerHelper = require('../../service/audio-book-player-helper')
const cardService = require('../card/card-service')
const cardHelper = require('../card/card-helper')

const getAllAudioBooks = async () => {
    const completeResult = await audioBookDbClient.findAllAudioBooks()

    return completeResult.map(doc => audioBookHelper.getAudioBookForResponse(doc))
}

const addAudioBook = async (audioBook) => {
    const result = await audioBookDbClient.insertAudioBook(audioBook)

    if (audioBook.card) {
        await cardService.removeUnusedCard(audioBook.card.id)
    }

    const mergedAudioBook = audioBookHelper.mergeAudioBook(result, audioBook)
    createMediaFolder(mergedAudioBook.id)

    return result
}

const updateAudioBook = async (id, audioBook) => {
    const result = await audioBookDbClient.findOneAudioBookById(id)

    if (result.card && audioBook.card && result.card.cardId !== audioBook.card.cardId) {
        await cardService.addUnusedCard(cardHelper.buildEmptyDbCard(result.card.cardId))
    }
    if (audioBook.card) {
        await cardService.removeUnusedCard(audioBook.card.id)
    }

    const mergedAudioBook = audioBookHelper.mergeAudioBook(audioBook, result)

    return audioBookDbClient.updateAudioBook(mergedAudioBook)
}

const deleteAudioBook = async (audioBookId) => {
    const audioBook = await audioBookDbClient.findOneAudioBookById(audioBookId)
    const result = await audioBookDbClient.deleteAudioBook(audioBook)
    deleteMediaFolder(audioBookId)

    if (audioBook.card) {
        await cardService.addUnusedCard(cardHelper.buildEmptyDbCard(audioBook.card.cardId))
    }

    return {
        audioBook: audioBook,
        dbResponse: result
    }
}

const uploadTracks = async (audioFiles, audioBookId) => {
    const audioBook = await audioBookDbClient.findOneAudioBookById(audioBookId)

    const trackLists = []

    let trackNumber = 1

    for (const audioFileName of Object.keys(audioFiles)) {
        const audioFile = audioFiles[audioFileName]
        const uploadPath = path.join(process.cwd(), 'media', audioBookId, audioFile.name)

        await audioFile.mv(uploadPath)

        const mediaInfo = await audioBookHelper.getMediaInfo(uploadPath)

        const trackList = audioBookHelper.buildNewAudioBookTrackList(trackNumber++, mediaInfo.title, mediaInfo.duration, audioFile.name)
        audioBook.trackList.push(trackList)
        trackLists.push(trackList)
    }

    await audioBookDbClient.updateAudioBook(audioBook)

    return trackLists
}

const deleteAllTracks = async (audioBookId) => {
    const audioBook = await audioBookDbClient.findOneAudioBookById(audioBookId)
    logger.info('Deleting Audio files for audio book: ' + audioBook.title)

    deleteMediaFolderContent(audioBook)
    audioBook.trackList = []

    const result = await audioBookDbClient.updateAudioBook(audioBook)

    return {
        audioBook: audioBook,
        dbResponse: result
    }
}

const createMediaFolder = (id) => {
    const mediaPath = getCompletePathToMediaFolder(id)
    fs.mkdirSync(mediaPath, { recursive: true })
}

const deleteMediaFolder = (id) => {
    const mediaPath = getCompletePathToMediaFolder(id)
    deleteMediaFolderRecursive(mediaPath)
}

const deleteMediaFolderContent = (audioBook) => {
    for (let track of audioBook.trackList) {
        const mediaPath = audioBookPlayerHelper.getCompletePathToAudioFile(audioBook, track)
        fs.unlinkSync(mediaPath)
    }
}

const deleteMediaFolderRecursive = (currentPath) => {
    if (fs.existsSync(currentPath)) {
        fs.readdirSync(currentPath).forEach(function (file) {
            const curPath = currentPath + '/' + file
            if (fs.lstatSync(curPath).isDirectory()) {
                deleteMediaFolderRecursive(curPath)
            } else {
                fs.unlinkSync(curPath)
            }
        })
        fs.rmdirSync(currentPath)
    }
}

const getCompletePathToMediaFolder = (id) => {
    return path.join(process.cwd(), 'media', id)
}

module.exports = {
    getAllAudioBooks,
    addAudioBook,
    updateAudioBook,
    deleteAudioBook,
    uploadTracks,
    deleteAllTracks
}
