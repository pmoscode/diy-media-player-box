const bodyParser = require('body-parser')
const logger = require('../logger')(module)
const express = require('express')
const helmet = require('helmet')
const listEndpoints = require('express-list-endpoints')
const fileUpload = require('express-fileupload')
const path = require('path')

exports.createAndInit = async () => {
    const app = express()
    await init(app)

    return app
}

const init = async (app) => {
    app.use(helmet())
    app.use(bodyParser.urlencoded({
        limit: '30mb',
        extended: true
    }))
    app.use(bodyParser.json())
    app.use(fileUpload({
        useTempFiles: true,
        tempFileDir: '/tmp/'
    }))
    // app.use(function (req, res, next) {
    //     res.header('Access-Control-Allow-Origin', '*')
    //     res.header('Access-Control-Allow-Methods', 'PATCH,POST,GET,DELETE')
    //     res.header('Access-Control-Allow-Headers', 'content-type')
    //     next()
    // })

    app.use(express.static(path.join(process.cwd(), 'ui/public')))

    require('../api/health/health-api').init(app)
    require('../api/swagger/swagger-api').init(app)
    require('../api/audio-book/audio-book-api').init(app)
    require('../api/card/card-api').init(app)

    logRegisteredRoutes(app)

    if (isRfidEnabled()) {
        logger.info('Rfid capability enabled')

        const rfid = require('../rfid/rfid-service')
        rfid.initRfid()
    } else {
        logger.info('Rfid capability NOT enabled')
    }
}

const logRegisteredRoutes = function (app) {
    function space (x) {
        let res = ''
        for (; x > 0; x--) {
            res += ' '
        }
        return res
    }

    const endpoints = listEndpoints(app) || []

    logger.info('#################################################################')
    logger.info('### DISPLAYING REGISTERED ROUTES:                             ###')
    logger.info('###                                                           ###')
    endpoints.forEach((r) => {
        r.methods.forEach((e) => {
            logger.info('### ' + e + space(8 - e.length) + r.path + space(50 - r.path.length) + '###')
        })
    })
    logger.info('###                                                           ###')
    logger.info('#################################################################')
}

const isRfidEnabled = () => {
    if (process.env.RFID_ENABLED === undefined) {
        return true
    }

    return JSON.parse(process.env.RFID_ENABLED)
}
