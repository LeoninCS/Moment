import { createRouter, createWebHistory } from 'vue-router'

import Home from '../views/Home.vue'
import CodeShare from '../views/CodeShare.vue'
import CodeView from '../views/CodeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/codeshare', name: 'codeshare', component: CodeShare },
    { path: '/code/:hash', name: 'codeview', component: CodeView, props: true },
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ],
  scrollBehavior() { return { top: 0 } }
})

export default router
