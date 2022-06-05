const musicMetadata = require('music-metadata')

const buildEmptyDbAudioBook = () => {
    return {
        id: '',
        title: '',
        lastPlayed: null,
        card: null,
        timesPlayed: 0,
        trackList: []
    }
}

const mergeAudioBook = (audioBookSource, audioBookTarget) => {
    return Object.assign(audioBookTarget, audioBookSource)
}

const getAudioBookFromRequest = (req) => {
    if (req && req.body) {
        return mergeAudioBook(req.body, buildEmptyDbAudioBook())
    }
}

const getAudioBookForResponse = (audioBook) => {
    if (audioBook) {
        if (audioBook._id) {
            audioBook.id = audioBook._id
        }
        delete audioBook.ok
        delete audioBook._id
        delete audioBook._rev
        delete audioBook.rev

        return audioBook
    }
}

const buildNewAudioBookTrackList = (track, title, length, filename) => {
    return {
        track: track,
        title: title,
        length: length,
        fileName: filename
    }
}

const getMediaInfo = async (file) => {
    const audioFileMetadata = await musicMetadata.parseFile(file)

    const duration = new Date(1000 * Math.round(audioFileMetadata.format.duration)).toISOString().substr(11, 8)
    const title = audioFileMetadata.common.title
    const artist = audioFileMetadata.common.artist

    return {
        duration: duration,
        title: title,
        artist: artist
    }
}

module.exports = {
    mergeAudioBook,
    getAudioBookFromRequest,
    getAudioBookForResponse,
    buildNewAudioBookTrackList,
    getMediaInfo
}
