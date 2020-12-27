require('dotenv').config()

const app = require('./app/server/app')
const initializer = require('./app/server/server-initializer')

app.createAndInit().then((application) => {
    initializer.startListening(application)
})
