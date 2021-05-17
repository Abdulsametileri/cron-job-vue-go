import Vue from 'vue'
import i18n from './plugins/i18n';

import {BootstrapVue, IconsPlugin} from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import App from './App.vue'
import router from './router'
import store from './store'

import messages from '@/mixins/messages'

Vue.use(BootstrapVue)
Vue.use(IconsPlugin)

i18n.locale = process.env.VUE_APP_LOCALE || 'tr';

Vue.config.productionTip = false

Vue.mixin({
    data() {
        return {
            appLocal: i18n.locale,
            isDevelopment: process.env.NODE_ENV === "development",
            repeatTypes: [
                {value: -1, text: this.$t('weekDays.default')},
                {value: 1, text: this.$t('weekDays.monday')},
                {value: 2, text: this.$t('weekDays.tuesday')},
                {value: 3, text: this.$t('weekDays.wednesday')},
                {value: 4, text: this.$t('weekDays.thursday')},
                {value: 5, text: this.$t('weekDays.friday')},
                {value: 6, text: this.$t('weekDays.saturday')},
                {value: 0, text: this.$t('weekDays.sunday')},
                {value: 7, text: this.$t('weekDays.all')},
            ],
            indexStrToWeekDay: {
                "-1": this.$t('weekDays.default'),
                "1": this.$t('weekDays.monday'),
            }
        };
    },
});

Vue.mixin(messages)

new Vue({
    router,
    store,
    i18n,
    render: h => h(App)
}).$mount('#app')
