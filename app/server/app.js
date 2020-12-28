const bodyParser = require('body-parser')
const logger = require('../helper/logger')(module)
const express = require('express')
const helmet = require('helmet')
const listEndpoints = require('express-list-endpoints')
const fileUpload = require('express-fileupload')
const environmentHelper = require('../helper/environment')

exports.createAndInit = () => {
    const app = express()
    const httpRequestSizeLimit = environmentHelper.getHttpRequestSizeLimit()

    app.use(helmet())
    app.use(bodyParser.urlencoded({
        limit: httpRequestSizeLimit,
        extended: true
    }))
    app.use(bodyParser.json())
    app.use(fileUpload({
        useTempFiles: true,
        tempFileDir: '/tmp/'
    }))
    if(environmentHelper.getAllowAllOrigin()) {
        app.use(function (req, res, next) {
            res.header('Access-Control-Allow-Origin', '*')
            res.header('Access-Control-Allow-Methods', 'PATCH,POST,GET,DELETE')
            res.header('Access-Control-Allow-Headers', 'content-type')
            next()
        })
    }

    app.use(express.static(environmentHelper.getStaticContentPath()))

    require('../api/health/health-api').init(app)
    require('../api/swagger/swagger-api').init(app)
    require('../api/audio-book/audio-book-api').init(app)
    require('../api/card/card-api').init(app)

    logRegisteredRoutes(app)

    return app
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
