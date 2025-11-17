import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import Home from '../views/Home.vue'
import CodeShare from '../views/CodeShare.vue'

import CodeShareView from '../views/CodeShareView.vue'
import WorkPlan from '../views/WorkPlan.vue'
import WorkPlanView from '../views/WorkPlanView.vue'
import FeatureC from '../components/FeatureC.vue'
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: Home,
  },
  {
    path: '/codeshare',
    name: 'codeshare',
    component: CodeShare,
  },
  {
    path: '/workplan',
    name: 'WorkPlan',
    component: WorkPlan,
  },
  {
    path: '/workplan/:hash',
    name: 'WorkPlanView',
    component: WorkPlanView,
    props: true,
  },
  {
    path: '/feature-c',
    name: 'feature-c',
    component: FeatureC,
  },
  {
    path: '/code/:hash',
    name: 'codeshare-view',
    component: CodeShareView,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
