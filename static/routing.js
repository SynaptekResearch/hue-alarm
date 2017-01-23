
var Vue = require('vue');
var VueRouter = require('vue-router');

import Home from './home.vue';
import Configuration from './config.vue';

Vue.use(VueRouter);

export var Router = new VueRouter({
  'routes': [
    { 'path': '/', 'component': Home },
    { 'path': '/config', 'component': Configuration }
  ]
});
