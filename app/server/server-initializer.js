const logger = require('../logger')(module)

const startListening = function (app) {
    process.on('uncaughtException', (ex) => {
        logger.error(ex)
    })
    app.listen(getPort(), logStartup)
}

const getPort = function () {
    return parseInt(process.env.PORT) || 8080
}

const logStartup = function () {
    logger.info('Application "DIY Media Player for children" is listening on port %s.', getPort())
}

module.exports = {
    startListening,
    getPort,
    logStartup
}
