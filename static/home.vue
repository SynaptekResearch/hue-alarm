<template>
  <div>
    <div class="ui one column row">
      <div class="ui raised segment column">
        <a class="ui green ribbon label">Status</a>
        <span></span>
        <p></p>
        
        <div v-if="status" style="font-size: 10pt">
          <span class="w2">Running: </span>
          <strong v-if="status.running" style="color:green">Yes</strong>
          <strong v-if="!status.running" style="color:red">No</strong><br/>
          <span class="w2">Armed: </span>
          <strong v-if="status.status.armed" style="color:green">Yes</strong>
          <strong v-if="!status.status.armed" style="color:red">No</strong><br/>
          <span class="w2">Last notified: </span><strong>{{status.status.lastNotified}}</strong><br/>
        </div>

        <div v-if="loading" class="ui loading button">...Loading...</div>
      </div>
      <div class="ui raised segment column">
        <a class="ui green ribbon label">Software information</a>
        <span></span>
        <p></p>

        <i>(c) Chris Polderman, GPL Licensed</i><br/>

        See <a href="https://www.github.com/cpo/hue-alarm">https://www.github.com/cpo/hue-alarm</a> for documentation.<br/><br/>
      </div>
    </div>
  </div>

</template>


<script>
var Vue = require('vue');
var api = require('./api');

var homeData = {
  status: null,
  loading: true
}

var timerId;

// Home screen
module.exports = Vue.component('home', {
  data: function() {
    return homeData;
  },
  beforeMount: function() {
    this.update();
    timerId = setInterval(this.update,5000);
  },
  beforeDestroy: function() {
    clearInterval(timerId);
  },
  methods: {
    update: function() {
        homeData.loading = true;
        api
          .get('/api/status')
          .then(function (response) {
            homeData.status = response.data;
            homeData.loading = false;
          })
          .catch(function (err) {
            console.log("Error accessing resource", err);
            homeData.loading = false;
          });
      }
  }
});

</script>