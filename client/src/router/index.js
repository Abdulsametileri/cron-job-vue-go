import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'AlarmList',
        component: () => import('@/views/AlarmList')
    },
    {
        path: '/create-alarm',
        name: 'AlarmCreate',
        component: () => import('@/views/AlarmCreate')
    }
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router
