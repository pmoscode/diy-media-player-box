const router = require('express').Router()
const swaggerValidator = require('../swagger/swagger-validator')
const apiErrorHandler = require('../api-error-handler')

exports.init = (app) => {
    swaggerValidator.get(router, '/health', healthStatus)

    app.use('/api', [router, apiErrorHandler])
}

const healthStatus = (req, res) => {
    return res.status(200).send({ status: 'OK' })
}
