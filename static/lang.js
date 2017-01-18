var Vue = require('vue');
var VueI18n = require('vue-i18n');

const locales = {
      en: {
        message: {
          name: 'HUE Alarm',
          menu_status: 'Status',
          menu_configuration: 'Configuration',
          config_activation: 'Activation',
          config_part_of_schedule_name: 'Trigger by part of schedule name',
          config_username: 'HUE username',
          config_primary_notification_settings: 'Primary notification settings',
          config_primary_notification_url: 'Notification URL (GET)',
          config_primary_notification_speed: 'Notification speed (seconds)',
          config_primary_notification_test_mode: 'Test mode',
          config_primary_notification_test: 'Test notification',
          config_secondary_notification_settings: 'Secundary notification settings',
          config_smtp_password: "SMTP password",
          config_smtp_server_name: "SMTP server name",
          config_smtp_port: "SMTP port",
          config_from: "FROM email address",
          config_to: "TO email address",
          config_message: "Message",
          enabled: "Enabled",
          config_admin_password: "Admin password",
          config_admin_username: "Admin username",
          config_administration: "Administration"
        }
      },
      nl: {
        message: {
          name: 'HUE Alarm',
          menu_status: 'Status',
          menu_configuration: 'Instellingen',
          config_activation: 'Activering',
          config_part_of_schedule_name: 'Activeer door woord in routine naam',
          config_username: 'HUE gebruikersnaam',
          config_primary_notification_settings: 'Hoofd notificatie instellingen',
          config_primary_notification_url: 'Notificatie URL (GET)',
          config_primary_notification_speed: 'Notification snelheid (seconden)',
          config_primary_notification_test_mode: 'Test modus',
          config_primary_notification_test: 'Test notificatie',
          config_secondary_notification_settings: 'Secundaire notificatie instellingen',
          config_smtp_password: 'SMTP wachtwoord',
          config_smtp_server_name: "SMTP server naam",
          config_smtp_port: "SMTP poort",
          config_from: "FROM email adres",
          config_to: "TO email adres",
          config_message: "Informatie",
          enabled: "Enabled",
          config_admin_password: "Admin wachtwoord",
          config_admin_username: "Admin gebruikersnaam",
          config_administration: "Beheer"

        }
      }
    };

Vue.use(VueI18n);
Vue.config.lang = 'en';

Object.keys(locales).forEach(function (lang) {
  Vue.locale(lang, locales[lang])
});

