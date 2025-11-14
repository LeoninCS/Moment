import { createRouter, createWebHistory } from 'vue-router'

import Home from '../views/Home.vue'
import CodeShare from '../views/CodeShare.vue'
import CodeView from '../views/CodeView.vue'
import EncryptandDecryptText from '../views/TextCrypto.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: Home },

    { path: '/codeshare', name: 'codeshare', component: CodeShare },

    { path: '/code/:hash', name: 'codeview', component: CodeView, props: true },

    { path: '/textcrypto', name: 'textcrypto', component: EncryptandDecryptText },


    // ⚠️ 404 通配路由一定要放在最后
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ],
  scrollBehavior() {
    return { top: 0 }
  }
})

export default router
