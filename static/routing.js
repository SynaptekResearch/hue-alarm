
var Vue = require('vue');
var VueRouter = require('vue-router');
var Home = require('./home.vue');
var Config = require('./config.vue');


Vue.use(VueRouter);

const routes = [
  {'path':'/', 'component': Home},
  {'path':'/config', 'component': Config}
];

module.exports = new VueRouter({
  'routes': routes 
});

