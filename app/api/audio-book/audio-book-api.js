const router = require('express').Router()
const swaggerValidator = require('../swagger/swagger-validator')
const apiErrorHandler = require('../api-error-handler')
const audioBookService = require('./audio-book-service')
const audioBookTestingService = require('./audio-book-testing-service')
const audioBookHelper = require('./audio-book-helper')

const wrap = fn => (...args) => fn(...args).catch(args[2])

exports.init = (app) => {
    swaggerValidator.get(router, '/audio-books', wrap(getAllAudioBooks))
    swaggerValidator.post(router, '/audio-books', wrap(addAudioBook))
    swaggerValidator.patch(router, '/audio-books/:id', wrap(updateAudioBook))
    swaggerValidator.delete(router, '/audio-books/:id', wrap(deleteAudioBook))
    swaggerValidator.post(router, '/audio-books/:id/tracks', wrap(uploadTracks))
    swaggerValidator.delete(router, '/audio-books/:id/tracks', wrap(deleteAllTracks))

    // Testing api
    swaggerValidator.post(router, '/audio-books/:id/track/:track/play', wrap(playTrack))
    swaggerValidator.post(router, '/audio-books/pause', wrap(pauseTrack))
    swaggerValidator.post(router, '/audio-books/stop', wrap(stopTrack))

    app.use('/api', [router, apiErrorHandler])
}

const getAllAudioBooks = async (req, res) => {
    let result = await audioBookService.getAllAudioBooks()

    if (!result) {
        result = []
    }

    return res.status(200).send(result)
}

const addAudioBook = async (req, res) => {
    const audioBook = audioBookHelper.getAudioBookFromRequest(req)

    const result = await audioBookService.addAudioBook(audioBook)

    if (result.ok) {
        const mergedAudioBook = audioBookHelper.mergeAudioBook(result, audioBook)

        return res.status(200).send(audioBookHelper.getAudioBookForResponse(mergedAudioBook))
    } else {
        return res.status(500).send('AudioBook not inserted!')
    }
}

const updateAudioBook = async (req, res) => {
    const audioBook = audioBookHelper.getAudioBookFromRequest(req)
    const id = req.params.id

    const result = await audioBookService.updateAudioBook(id, audioBook)

    if (result.ok) {
        const mergedAudioBook = audioBookHelper.mergeAudioBook(result, audioBook)

        return res.status(200).send(audioBookHelper.getAudioBookForResponse(mergedAudioBook))
    } else {
        return res.status(500).send('AudioBook not updated!')
    }
}

const deleteAudioBook = async (req, res) => {
    const id = req.params.id

    const result = await audioBookService.deleteAudioBook(id)

    if (result.dbResponse.ok) {
        return res.status(200).send(audioBookHelper.getAudioBookForResponse(result.audioBook))
    } else {
        return res.status(500).send('AudioBook not deleted!')
    }
}

const uploadTracks = async (req, res) => {
    const files = req.files
    const audioBookId = req.params.id

    if (Object.keys(files).length === 0) {
        return res.status(400).send('No files were uploaded.')
    }

    try {
        const track = await audioBookService.uploadTracks(files, audioBookId)

        return res.status(200).send(track)
    } catch (e) {
        return res.status(500).send(e)
    }
}

const deleteAllTracks = async (req, res) => {
    const audioBookId = req.params.id

    const result = await audioBookService.deleteAllTracks(audioBookId)

    if (result.dbResponse.ok) {
        return res.status(200).send(audioBookHelper.getAudioBookForResponse(result.audioBook))
    } else {
        return res.status(500).send('AudioBook tracks not deleted!')
    }
}

const playTrack = async (req, res) => {
    const audioBookId = req.params.id
    const trackNumber = req.params.track

    try {
        await audioBookTestingService.play(audioBookId, trackNumber)

        return res.status(200).send({ status: 'playing' })
    } catch (e) {
        return res.status(500).send(e)
    }
}

const pauseTrack = async (req, res) => {
    try {
        audioBookTestingService.pause()

        return res.status(200).send({ status: 'play state toggled' })
    } catch (e) {
        return res.status(400).send({ status: 'no track is playing...' })
    }
}

const stopTrack = async (req, res) => {
    try {
        audioBookTestingService.stop()

        return res.status(200).send({ status: 'stopped' })
    } catch (e) {
        return res.status(400).send({ status: 'no track is playing...' })
    }
}
