// TODO: This service must be implemented in GO

const logger = require('../helper/logger')(module)
const helper = require('./audio-book-player-helper')

let lastPlayedUid = 0
let checkProgressInterval
let currentAudioBook
let currentAudioTrack = 0

const nfcNoCardDetected = () => {
    logger.debug('player toggle')
    try {
        clearInterval(checkProgressInterval)
        audioPlayer.togglePause()
    } catch (e) {
        logger.error(e)
    }
}

const nfcCardDetected = (uid) => {
    logger.info('LastPlayedUid: ' + lastPlayedUid)
    if (lastPlayedUid === uid) {
        logger.debug('player toggle')
        try {
            audioPlayer.togglePause()
        } catch (e) {
            if (currentAudioBook.trackList.length > currentAudioTrack) {
                playTrack()
            } else {
                logger.error(e)
            }
        }
    } else {
        logger.debug('player start')
        clearInterval(checkProgressInterval)
        helper.getFilePathForUid(uid).then((result) => {
            if (result.filePath) {
                currentAudioBook = result.audioBook
                currentAudioTrack = 0
                playAudioFile(result.filePath, result.uid)
                helper.updateMetadata(currentAudioBook).then(result => {
                    logger.info('Update metadata for uid: ' + uid + ' and audioBook: ' + currentAudioBook.title + ' is ' + result.ok)
                })
            } else {
                try {
                    audioPlayer.stop()
                } catch (e) {
                    logger.debug('Panic audio stop...')
                }
            }
        }).catch()
    }
}

const playAudioFile = (filePath, uid) => {
    audioPlayer.play(filePath)
        .then(() => {
            logger.info('Playing audio file "' + filePath + '"')
            lastPlayedUid = uid
            monitorProgress()
        })
        .catch(e => {
            logger.error(e)
        })
}

const playNextTrack = () => {
    ++currentAudioTrack

    playTrack()
}

const playTrack = () => {
    if (currentAudioBook) {
        const track = currentAudioBook.trackList[currentAudioTrack]

        if (track) {
            const completePath = helper.getCompletePathToAudioFile(currentAudioBook, track)
            playAudioFile(completePath, lastPlayedUid)
        } else {
            logger.info('No more tracks found. Finished...')
            lastPlayedUid = 0
            currentAudioTrack = undefined
            currentAudioBook = undefined
        }
    }
}

module.exports = {
    nfcNoCardDetected,
    nfcCardDetected
}
