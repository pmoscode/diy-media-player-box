const getValueFromEnvironment = (key, defaultValue) => {
    const value = process.env[key]

    return value ? value : defaultValue
}

const parseBool = (bool) => {
    if(typeof bool === 'boolean') {
        return bool
    }

    return bool.toLowerCase() === 'true'
}

const parseInteger = (int) => {
    if(typeof int === 'number') {
        return int
    }

    return parseInt(int)
}

module.exports = {
    getValueFromEnvironment,
    parseBool,
    parseInteger
}
