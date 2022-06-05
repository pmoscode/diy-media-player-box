const logger = require('../helper/logger')(module)
const environmentHelper = require('../helper/environment')

const startListening = function (app) {
    process.on('uncaughtException', (ex) => {
        logger.error(ex.toString())
    })
    app.listen(environmentHelper.getAppServerPort(), logStartup)
}

const logStartup = function () {
    logger.info('Application "DIY Media Player for children" is listening on port %s.', environmentHelper.getAppServerPort())
}

module.exports = {
    startListening
}
