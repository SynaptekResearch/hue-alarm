
var Vue = require('vue');

var $ = require('jquery');
window.jQuery=$;
var semantic = require('semantic-ui/dist/semantic.min');

require('./lang');
var router = require('./routing');

// Main application
var app = new Vue({
  router,
  el: '#app'
});

$(document).ready(function() {
  // fix main menu to page on passing
  $('.main.menu').visibility({
    type: 'fixed'
  });
  $('.overlay').visibility({
    type: 'fixed',
    offset: 80
  });

  // lazy load images
  $('.image').visibility({
    type: 'image',
    transition: 'vertical flip in',
    duration: 500
  });

  // show dropdown on hover
  $('.main.menu .ui.dropdown').dropdown({
    on: 'hover'
  });
});