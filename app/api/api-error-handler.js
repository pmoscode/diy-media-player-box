const logger = require('../helper/logger')(module)

function handle (err, req, res, next) {
    logger.error('handling error: %s', err)

    let statusCode = err.statusCode ? err.statusCode : err.status ? err.status : 500

    if (err.name === 'ValidationError') {
        statusCode = 400
    }

    res.status(statusCode).json({ statusCode: statusCode, message: err.message })
}

module.exports = handle
