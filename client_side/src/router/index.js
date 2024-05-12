
import { createRouter, createWebHistory } from 'vue-router'
const routerHistory = createWebHistory()
const router = createRouter({
    history: routerHistory,
    routes: [
        {
            path: '/index.html',
            redirect: '/'
        },
        {
            path: '/',
            redirect: '/index'
        },
        {
            path: '/index',
            redirect: '/index/Overview',
            component: () => import('@/components/Index.vue'),
            children:[
                {path: '/index/Overview',component: () => import('@/components/Overview/Overview.vue')},
                {path: '/index/CFG',component: () => import('@/components/CFG/CFG.vue')},
                {path: '/index/Seeds',component: () => import('@/components/Seeds/Seeds.vue')},
                {path: '/index/Terminal-main',component: () => import('@/components/Terminal/Terminal-main.vue')},
                {path: '/index/Terminal-1',component: () => import('@/components/Terminal/Terminal-1.vue')}
            ]
        },

    ]
})
export default router
