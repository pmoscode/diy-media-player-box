const path = require('path')
const audioPlayer = require('../../audio-player/audio-player-controller')
const audioBookDbClient = require('../../service/audio-book-db-client')
const logger = require('../../helper/logger')(module)

const play = async (audioBookId, trackNumber) => {
    const audioBook = await audioBookDbClient.findOneAudioBookById(audioBookId)
    const trackList = audioBook.trackList[trackNumber - 1]

    const completePath = path.join(process.cwd(), 'media', audioBookId, trackList.fileName)
    logger.info('Playing audio "%s"', completePath)

    await audioPlayer.play(completePath)
}

const pause = () => {
    audioPlayer.togglePause()
}

const stop = () => {
    audioPlayer.stop()
}

module.exports = {
    play,
    pause,
    stop
}
