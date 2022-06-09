const path = require('path')
const fs = require('fs')
const axios = require('axios')
const logger = require('../../helper/logger')(module)
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

const playAudioBook = async (audioBookId, trackNumber) => {
    const audioBook = await audioBookDbClient.findOneAudioBookById(audioBookId)
    const trackList = audioBook.trackList[trackNumber - 1]

    const completePath = path.join(process.cwd(), 'media', audioBookId, trackList.fileName)
    logger.info('Playing audio "%s"', completePath)

    const requestBody = {
        uid: audioBookId,
        trackList: [completePath]
    }

    axios.post('http://localhost:8080/audio', requestBody)
        .then((res) => {
            console.log(`Status: ${res.status}`)
            console.log('Student Info: ', res.data)
        }).catch((err) => {
            console.error(err)
        })
}

const playAudioBookFromUid = async (uid) => {
    const tryAudioBook = await audioBookPlayerHelper.getFilePathForUid(uid)
    if (tryAudioBook.audioBook) {
        const trackList = tryAudioBook.audioBook.trackList

        const trackListPath = []

        for (const track of trackList) {
            const completePath = path.join(process.cwd(), 'media', uid, track.fileName)
            // logger.info('Playing audio "%s"', completePath)
            trackListPath.push(completePath)
        }

        const requestBody = {
            uid: uid,
            trackList: trackListPath
        }

        axios.post('http://localhost:8080/audio', requestBody)
            .then((res) => {
                console.log(`Status: ${res.status}`)
                console.log('Response: ', res.data)
            }).catch((err) => {
                console.error(err)
            })
    } else {
        console.log('Seems the card used is not known...: ', uid)
    }
}

const pauseAudioBook = () => {
    axios.patch('http://localhost:8080/audio')
        .then((res) => {
            console.log(`Status: ${res.status}`)
            console.log('Response: ', res.data)
        }).catch((err) => {
            console.error(err)
        })
}

const stopAudioBook = () => {
    axios.delete('http://localhost:8080/audio')
        .then((res) => {
            console.log(`Status: ${res.status}`)
            console.log('Response: ', res.data)
        }).catch((err) => {
            console.error(err)
        })
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
    for (const track of audioBook.trackList) {
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
    deleteAllTracks,
    playAudioBook,
    playAudioBookFromUid,
    pauseAudioBook,
    stopAudioBook
}
