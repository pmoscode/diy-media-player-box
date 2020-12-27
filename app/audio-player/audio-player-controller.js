const mpg = require('mpg123')
const player = new mpg.MpgPlayer()
const fs = require('fs').promises

const playerGetProgress = () => {
    return new Promise((resolve, reject) => {
        try {
            player.getProgress((p, s, l) => {
                resolve({ p, s, l })
            })
        } catch (e) {
            reject(e)
        }
    })
}

const play = async (filePath) => {
    try {
        await fs.access(filePath)

        player.play(filePath)
    } catch (error) {
        throw new Error('File does not exists!')
    }
}

const togglePause = () => {
    if (!player.track) {
        throw new Error('Currently nothing is played!')
    }

    player.pause()
}

const stop = () => {
    if (!player.track) {
        throw new Error('Currently nothing is played!')
    }

    player.stop()
}

const getProgress = () => {
    if (!player.track) {
        throw new Error('Currently nothing is played!')
    }

    return playerGetProgress()
}

module.exports = {
    play,
    togglePause,
    stop,
    getProgress
}
