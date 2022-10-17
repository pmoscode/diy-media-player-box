const jsYaml = require('js-yaml')
const path = require('path')
const fs = require('fs')

exports.getSwaggerDocument = () => {
    return jsYaml.safeLoad(fs.readFileSync(path.resolve(__dirname, './swagger.yaml'), 'utf-8'))
}
