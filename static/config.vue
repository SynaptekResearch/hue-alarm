
<template>
  <div>
    <div class="one column row">
      <form class="ui form" v-if="config && config.statusMessages">
        <div class="ui two column row">
          <div class="ui raised segment column">
            <a class="ui green ribbon label">{{ $t("message.config_activation") }}</a>
            <span></span>
            <p></p>

            <div class="field">
              <label>
                {{ $t("message.config_part_of_schedule_name") }}
              </label>
              <input type="text" v-model="config.triggerOnSchedulePart" placeholder="Keyword in schedule name">
            </div>

            <div class="field">
              <label>
                {{ $t("message.config_username") }}
              </label>
              <input type="text" v-model="config.userName" placeholder="HUE username">
            </div>

          </div>
        </div>

        <div class="ui raised segment">
          <a class="ui green ribbon label">{{ $t("message.config_primary_notification_settings") }}</a>
          <span></span>
          <p></p>

          <div class="field">
            <label>
              {{ $t("message.config_primary_notification_url") }}
            </label>
            <input v-model="config.notificationURL" placeholder="Notification GET URL">
          </div>

          <div class="field">
            <label>
              {{ $t("message.config_primary_notification_speed") }}
            </label>
             <input v-model.number="config.notificationDelaySeconds" placeholder="Notification minimum interleave in seconds">
          </div>

          <div class="field">
            <label>
              {{ $t("message.config_primary_notification_test_mode") }}
            </label>
            <input type="checkbox" v-model="config.testMode" >
          </div>

          <div class="field">
            <button class="ui secondary button" v-on:click.stop="testnotification">{{ $t("message.config_primary_notification_test") }}</button>
          </div>

        </div>

        <div class="ui raised segment">
          <a class="ui green ribbon label">{{$t('message.config_secondary_notification_settings')}}</a>
          <span></span>
          <p></p>

          <div class="field">
            <label>
              {{$t('message.enabled')}}
            </label>
            <input type="checkbox" v-model="config.statusMessages.enabled" >
          </div>

          <div class="field">
            <label>
              {{$t('message.config_from')}}
            </label>
            <input v-model="config.statusMessages.from" placeholder="Enter FROM mail address">
          </div>

          <div class="field">
            <label>
              {{$t('message.config_to')}}
            </label>
            <input v-model="config.statusMessages.to" placeholder="Enter TO mail address">
          </div>

          <div class="field">
            <label>
              {{$t('message.config_smtp_server_name')}} 
            </label>
            <input v-model="config.statusMessages.smtpServer" placeholder="Enter SMTP server name">
          </div>

          <div class="field">
            <label>
              {{$t('message.config_smtp_port')}} 
            </label>
             <input v-model.number="config.statusMessages.smtpPort" placeholder="Enter SMTP port number">
          </div>

          <div class="field">
            <label>
              {{$t('message.config_smtp_password')}} 
            </label>
             <input v-model="config.statusMessages.password" placeholder="Enter SMTP password">
          </div>

          <button class="ui primary button" v-on:click.stop="saveSettings">Save</button>
          
          <div class="ui info message" v-if="message">
            <div class="header">
              {{$t('message.config_message')}} 
            </div>
            <p>{{message}}</p>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
var Vue = require('vue');
var api = require('./api');

var data = {
  config: {
    statusMessages: {}
  },
  message: ''
};

// Configuration screen
module.exports = Vue.component('configuration', {
  data: function() {
    return data;
  },
  created: function () {
    api
      .get('/api/config')
      .then(function (response) {
        data.config = response.data;
      })
      .catch(function (err) {
        console.log("Error accessing resource", err);
      });
  },
  methods: {
    testnotification: function() {
      api
        .post('/api/test-notify', {
          URL: data.config.notificationURL
        })
        .then(function (response) {
          data.message = response.data;
        });
    },
    saveSettings: function() {
      api
        .post('/api/config', data.config)
        .then(function (response) {
          data.message = response.data;
        });
    }
  }
});
</script>