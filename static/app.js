
import Vue from 'vue';

// dude, serious?
var $ = require('jquery');
window.jQuery=$;
// ...because aliens
var semantic = require('semantic-ui/dist/semantic.min');

import { activateLanguage } from './lang';
import { Router } from './routing';

console.log("Router is %o", Router);

activateLanguage();

// Main application
var app = new Vue({
  router: Router,
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