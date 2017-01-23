var axios = require('axios');


export var API = axios.create({
  timeout: 5000,
  headers: {
    'Pragma': 'no-cache',
    'Cache-Control': 'no-cache'
  }
});
