require('dotenv').config()

const app = require('./app/server/app')
const initializer = require('./app/server/server-initializer')
const environmentHelper = require('./app/helper/environment')
const logger = require('./app/helper/logger')(module)

// Express server init
const appInstance = app.createAndInit()
initializer.startListening(appInstance)

// RFID init
if (environmentHelper.isRfidEnabled()) {
    logger.info('RFID capability enabled')

    const rfid = require('./app/rfid/rfid-service')
    rfid.initRfid()
} else {
    logger.info('RFID capability NOT enabled')
}
