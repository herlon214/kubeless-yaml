const axios = require('axios')

module.exports = {
    handler: async (event, context) => {
        console.log(event, context)
        const resp = await axios.get('https://github.com')

        console.log(resp)
        return event.data
    }
}