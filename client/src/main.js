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

i18n.locale = 'tr';

Vue.config.productionTip = false

Vue.mixin({
    data() {
        return {
            isDevelopment: process.env.NODE_ENV === "development",
            repeatTypes: [
                {value: -1, text: i18n.tc('weekDays.default')},
                {value: 1, text: i18n.tc('weekDays.monday')},
                {value: 2, text: i18n.tc('weekDays.tuesday')},
                {value: 3, text: i18n.tc('weekDays.wednesday')},
                {value: 4, text: i18n.tc('weekDays.thursday')},
                {value: 5, text: i18n.tc('weekDays.friday')},
                {value: 6, text: i18n.tc('weekDays.saturday')},
                {value: 0, text: i18n.tc('weekDays.sunday')},
                {value: 7, text: i18n.tc('weekDays.all')},
            ],
            indexStrToWeekDay: {
                "-1": i18n.tc('weekDays.default'),
                "1": i18n.tc('weekDays.monday'),
                "2": i18n.tc('weekDays.tuesday'),
                "3": i18n.tc('weekDays.wednesday'),
                "4": i18n.tc('weekDays.thursday'),
                "5": i18n.tc('weekDays.friday'),
                "6": i18n.tc('weekDays.saturday'),
                "0": i18n.tc('weekDays.sunday'),
                "7": i18n.tc('weekDays.all'),
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
