var axios = require('axios');


module.exports = axios.create({
  timeout: 5000,
  headers: {
    'Pragma': 'no-cache',
    'Cache-Control': 'no-cache'
  }
});