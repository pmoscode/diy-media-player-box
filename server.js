require('dotenv').config()

const app = require('./app/server/app')
const initializer = require('./app/server/server-initializer')

// Express server init
const appInstance = app.createAndInit()
initializer.startListening(appInstance)
