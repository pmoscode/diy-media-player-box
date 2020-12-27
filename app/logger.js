const path = require('path')
const { grey, italic } = require('colors')
const { createLogger, format, transports } = require('winston')
const { combine, label, printf, colorize, splat } = format

const customFormat = printf(function ({ level, message, label }) {
    return `${level}: ${message} ${label}`
})

const getLabel = function (callingModule) {
    if (typeof callingModule === 'string') {
        return 'diy-audio-book-player/' + callingModule
    }

    const parts = callingModule.filename.split(path.sep)

    return path.join('diy-audio-book-player', parts[parts.length - 2], parts.pop())
}

module.exports = function (callingModule) {
    return createLogger({
        level: 'debug',
        format: combine(
            label({ label: grey(italic('- ' + getLabel(callingModule))) }),
            colorize(),
            splat(),
            customFormat
        ),
        transports: [new transports.Console(), new transports.File({
            filename: path.join(process.cwd(), 'logs', 'log_file.log'),
            tailable: true,
            maxsize: 200000
        })]
    })
}
