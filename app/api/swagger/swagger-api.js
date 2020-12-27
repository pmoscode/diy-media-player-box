const swaggerUi = require('swagger-ui-express')
const swaggerDocument = require('./swagger').getSwaggerDocument()

exports.init = (app) => {
    const swagger = swaggerUi.generateHTML(swaggerDocument)

    app.use('/api/swagger', swaggerUi.serveFiles(swaggerDocument))
    app.get('/api/swagger', (req, res) => res.send(swagger))
}
