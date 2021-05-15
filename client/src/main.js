import Vue from 'vue'
import i18n from './plugins/i18n';

import {BootstrapVue, IconsPlugin} from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import App from './App.vue'
import router from './router'
import store from './store'

Vue.use(BootstrapVue)
Vue.use(IconsPlugin)

i18n.locale = process.env.VUE_APP_LOCALE || 'tr';

Vue.config.productionTip = false

Vue.mixin({
    data() {
        return {
            appLocal: i18n.locale,
            isDevelopment: process.env.NODE_ENV === "development"
        };
    },
});

new Vue({
    router,
    store,
    i18n,
    render: h => h(App)
}).$mount('#app')
