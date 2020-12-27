const router = require('express').Router()
const swaggerValidator = require('../swagger/swagger-validator')
const apiErrorHandler = require('../api-error-handler')
const cardService = require('./card-service')

const wrap = fn => (...args) => fn(...args).catch(args[2])

exports.init = (app) => {
    swaggerValidator.get(router, '/cards/unassigned', wrap(getAllUnassignedCards))

    app.use('/api', [router, apiErrorHandler])
}

const getAllUnassignedCards = async (req, res) => {
    let result = await cardService.getAllUnassignedCards()

    if (!result) {
        result = []
    }

    return res.status(200).send(result)
}
