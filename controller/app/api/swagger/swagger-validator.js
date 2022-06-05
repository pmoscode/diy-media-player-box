const { OpenApiValidator } = require('express-openapi-validate')
const swaggerDocument = require('./swagger').getSwaggerDocument()

const validator = new OpenApiValidator(swaggerDocument, { ajvOptions: { allErrors: true } })

const _get = (app, path, callbacks) => {
    addRoute('get', app, path, callbacks)
}

const _post = (app, path, callbacks) => {
    addRoute('post', app, path, callbacks)
}

const _put = (app, path, callbacks) => {
    addRoute('put', app, path, callbacks)
}

const _patch = (app, path, callbacks) => {
    addRoute('patch', app, path, callbacks)
}

const _delete = (app, path, callbacks) => {
    addRoute('delete', app, path, callbacks)
}

function addRoute (method, app, path, callbacks) {
    // logger.info('adding validated route: %s %s', method, path)

    const openApiPath = convertToOpenApiPath(path)

    // logger.debug('calculated open api path=%s', openApiPath)

    if (!Array.isArray(callbacks)) {
        callbacks = [callbacks]
    }

    switch (method) {
    case 'get':
        app.get(path, [validator.validate(method, openApiPath)].concat(callbacks))
        break
    case 'post':
        app.post(path, [validator.validate(method, openApiPath)].concat(callbacks))
        break
    case 'put':
        app.put(path, [validator.validate(method, openApiPath)].concat(callbacks))
        break
    case 'patch':
        app.patch(path, [validator.validate(method, openApiPath)].concat(callbacks))
        break
    case 'delete':
        app.delete(path, [validator.validate(method, openApiPath)].concat(callbacks))
        break
    default:
        throw new Error('unknown method defined')
    }
}

function convertToOpenApiPath (path) {
    let openApiPath = ''

    path.split('/').map(name => {
        if (name.startsWith(':')) {
            return '/{' + name.substring(1) + '}'
        }

        return '/' + name
    }).forEach((name) => {
        openApiPath += name
    })

    return openApiPath.replace('//', '/')
}

module.exports = {
    get: _get,
    post: _post,
    put: _put,
    patch: _patch,
    delete: _delete
}
