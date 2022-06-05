const helper = require('./generic')
const path = require('path')

const getStaticContentPath = () => {
    return helper.getValueFromEnvironment('UI_CONTENT_PATH', path.join(process.cwd(), 'ui/public'))
}

const getHttpRequestSizeLimit = () => {
    return helper.getValueFromEnvironment('HTTP_REQUEST_SIZE_LIMIT', '30mb')
}

const getAllowAllOrigin = () => {
    return helper.parseBool(helper.getValueFromEnvironment('ALLOW_ALL_ORIGIN', false))
}

const getAppServerPort = () => {
    return helper.parseInteger(helper.getValueFromEnvironment('PORT', 8080))
}

module.exports = {
    getStaticContentPath,
    getHttpRequestSizeLimit,
    getAllowAllOrigin,
    getAppServerPort
}
