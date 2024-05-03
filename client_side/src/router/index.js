
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
                {path: '/index/Overview',component: () => import('@/components/body/Overview.vue')},
                {path: '/index/CFG',component: () => import('@/components/body/CFG.vue')},
                {path: '/index/Seeds',component: () => import('@/components/body/Seeds.vue')},
                {path: '/index/Terminal-main',component: () => import('@/components/body/Terminal-main.vue')},
                {path: '/index/Terminal-1',component: () => import('@/components/body/Terminal-1.vue')}
            ]
        },

    ]
})
export default router
